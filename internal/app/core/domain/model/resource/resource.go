package resource

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// warsawTZ is used to interpret HH:mm availability windows consistently with the web UI default.
var warsawTZ = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return time.UTC
	}
	return loc
}()

type Resource struct {
	id                    uuid.UUID
	companyID             uuid.UUID
	name                  string
	resourceType          ResourceType
	companyResourceTypeID *uuid.UUID
	customTypeName        *string // filled on read/list when JOINed; not authoritative for writes
	customTypeIconKey     *string
	description           string
	location              string
	maxReservationMinutes *int
	availableFrom         *string // HH:mm
	availableTo           *string // HH:mm
	status                ResourceStatus
	unavailableFrom       *time.Time
	unavailableTo         *time.Time
	unavailableReason     *string
	createdAt             time.Time
	updatedAt             time.Time
}

func NewResource(
	companyID uuid.UUID,
	name string,
	resType ResourceType,
	companyResourceTypeID *uuid.UUID,
	description string,
	location string,
	maxReservationMinutes *int,
	availableFrom *string,
	availableTo *string,
) (*Resource, error) {
	if name == "" {
		return nil, errors.New("resource name cannot be empty")
	}
	if resType == TypeCustom {
		if companyResourceTypeID == nil {
			return nil, errors.New("custom resource requires company_resource_type_id")
		}
	} else if companyResourceTypeID != nil {
		return nil, errors.New("company_resource_type_id is only allowed when type is custom")
	}
	if !IsBuiltInType(resType) && resType != TypeCustom {
		return nil, fmt.Errorf("unknown resource type: %s", resType)
	}

	return &Resource{
		id:                    uuid.New(),
		companyID:             companyID,
		name:                  name,
		resourceType:          resType,
		companyResourceTypeID: companyResourceTypeID,
		description:           description,
		location:              location,
		maxReservationMinutes: maxReservationMinutes,
		availableFrom:         availableFrom,
		availableTo:           availableTo,
		status:                StatusActive,
		createdAt:             time.Now().UTC(),
		updatedAt:             time.Now().UTC(),
	}, nil
}

func (r *Resource) MarkInService(reason string, from time.Time, to *time.Time) {
	r.status = StatusInService
	r.unavailableFrom = &from
	r.unavailableTo = to
	r.unavailableReason = &reason
	r.updatedAt = time.Now().UTC()
}

func (r *Resource) Restore() {
	r.status = StatusActive
	r.unavailableFrom = nil
	r.unavailableTo = nil
	r.unavailableReason = nil
	r.updatedAt = time.Now().UTC()
}

func (r *Resource) Deactivate() {
	r.status = StatusInactive
	r.updatedAt = time.Now().UTC()
}

// IsBookable is true when the resource can accept at least one new booking "right now"
// (e.g. UI: show booking button). Future maintenance does not block until it starts.
func (r *Resource) IsBookable() bool {
	switch r.status {
	case StatusInactive:
		return false
	case StatusActive:
		return true
	case StatusInService:
		now := time.Now().UTC()
		if r.unavailableFrom != nil && now.Before(*r.unavailableFrom) {
			return true
		}
		if r.unavailableTo != nil && !now.Before(*r.unavailableTo) {
			return true
		}
		return false
	default:
		return false
	}
}

// IsBookableForInterval is false when [start, end) overlaps scheduled in-service maintenance
// or falls outside configured daily HH:mm availability (Europe/Warsaw wall clock).
func (r *Resource) IsBookableForInterval(start, end time.Time) bool {
	switch r.status {
	case StatusInactive:
		return false
	case StatusActive:
		return r.intervalFitsDailyAvailability(start, end)
	case StatusInService:
		if r.unavailableFrom == nil {
			return false
		}
		if r.reservationOverlapsInService(start, end) {
			return false
		}
		return r.intervalFitsDailyAvailability(start, end)
	default:
		return false
	}
}

func parseHHMMToMinutes(s string) (int, bool) {
	var h, m int
	n, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	if err != nil || n != 2 {
		return 0, false
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, false
	}
	return h*60 + m, true
}

func truncateToMidnightInLoc(t time.Time, loc *time.Location) time.Time {
	y, mo, d := t.In(loc).Date()
	return time.Date(y, mo, d, 0, 0, 0, 0, loc)
}

// intervalFitsDailyAvailability verifies each Warsaw-local calendar slice of [start, end)
// lies inside [availableFrom, availableTo) expressed as HH:mm wall clocks when both pointers are set.
func (r *Resource) intervalFitsDailyAvailability(start, end time.Time) bool {
	if r.availableFrom == nil || r.availableTo == nil {
		return true
	}
	fromMin, okFrom := parseHHMMToMinutes(*r.availableFrom)
	toMin, okTo := parseHHMMToMinutes(*r.availableTo)
	if !okFrom || !okTo || fromMin >= toMin {
		return true // misconfigured availability does not hard-block bookings
	}
	if !start.Before(end) {
		return false
	}

	loc := warsawTZ
	open := start
	for open.Before(end) {
		dayMidnight := truncateToMidnightInLoc(open, loc)
		nextMidnight := dayMidnight.AddDate(0, 0, 1)

		segBegin := maxTime(open, dayMidnight)
		segEnd := minTime(end, nextMidnight)

		if !segBegin.Before(segEnd) {
			open = open.Add(time.Minute)
			continue
		}

		sh, sm, _ := segBegin.In(loc).Clock()
		eh, em, _ := segEnd.In(loc).Clock()
		startM := sh*60 + sm
		endM := eh*60 + em

		if startM < fromMin || endM > toMin {
			return false
		}

		open = segEnd
	}
	return true
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func halfOpenOverlap(s, e, a, b time.Time) bool {
	return s.Before(b) && e.After(a)
}

func (r *Resource) reservationOverlapsInService(start, end time.Time) bool {
	uf := *r.unavailableFrom
	if r.unavailableTo == nil {
		return end.After(uf)
	}
	ut := *r.unavailableTo
	return halfOpenOverlap(start, end, uf, ut)
}

func (r *Resource) ID() uuid.UUID                     { return r.id }
func (r *Resource) CompanyID() uuid.UUID              { return r.companyID }
func (r *Resource) Name() string                      { return r.name }
func (r *Resource) Type() ResourceType                { return r.resourceType }
func (r *Resource) CompanyResourceTypeID() *uuid.UUID { return r.companyResourceTypeID }
func (r *Resource) CustomTypeName() *string           { return r.customTypeName }
func (r *Resource) CustomTypeIconKey() *string        { return r.customTypeIconKey }
func (r *Resource) Description() string               { return r.description }
func (r *Resource) Location() string                  { return r.location }
func (r *Resource) MaxReservationMinutes() *int       { return r.maxReservationMinutes }
func (r *Resource) AvailableFrom() *string            { return r.availableFrom }
func (r *Resource) AvailableTo() *string              { return r.availableTo }
func (r *Resource) Status() ResourceStatus            { return r.status }
func (r *Resource) UnavailableFrom() *time.Time       { return r.unavailableFrom }
func (r *Resource) UnavailableTo() *time.Time         { return r.unavailableTo }
func (r *Resource) UnavailableReason() *string        { return r.unavailableReason }
func (r *Resource) CreatedAt() time.Time              { return r.createdAt }
func (r *Resource) UpdatedAt() time.Time              { return r.updatedAt }

type Props struct {
	ID                    uuid.UUID
	CompanyID             uuid.UUID
	Name                  string
	ResourceType          string
	CompanyResourceTypeID *uuid.UUID
	CustomTypeName        *string
	CustomTypeIconKey     *string
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
	Status                string
	UnavailableFrom       *time.Time
	UnavailableTo         *time.Time
	UnavailableReason     *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func Reconstitute(p *Props) *Resource {
	return &Resource{
		id:                    p.ID,
		companyID:             p.CompanyID,
		name:                  p.Name,
		resourceType:          ResourceType(p.ResourceType),
		companyResourceTypeID: p.CompanyResourceTypeID,
		customTypeName:        p.CustomTypeName,
		customTypeIconKey:     p.CustomTypeIconKey,
		description:           p.Description,
		location:              p.Location,
		maxReservationMinutes: p.MaxReservationMinutes,
		availableFrom:         p.AvailableFrom,
		availableTo:           p.AvailableTo,
		status:                ResourceStatus(p.Status),
		unavailableFrom:       p.UnavailableFrom,
		unavailableTo:         p.UnavailableTo,
		unavailableReason:     p.UnavailableReason,
		createdAt:             p.CreatedAt,
		updatedAt:             p.UpdatedAt,
	}
}
