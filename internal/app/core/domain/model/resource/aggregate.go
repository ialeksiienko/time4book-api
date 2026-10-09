package resource

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	id                    uuid.UUID
	companyID             uuid.UUID
	name                  string
	resourceType          ResourceType
	companyResourceTypeID *uuid.UUID
	description           string
	location              string
	maxReservationMinutes *int
	availableFrom         *string
	availableTo           *string
	status                ResourceStatus
	unavailableFrom       *time.Time
	unavailableTo         *time.Time
	unavailableReason     *string
	createdAt             time.Time
	updatedAt             time.Time
}

func NewResource(companyID uuid.UUID, details Details) (*Resource, error) {
	if err := validateDetails(details); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Resource{
		id:                    uuid.New(),
		companyID:             companyID,
		name:                  details.Name,
		resourceType:          details.Type,
		companyResourceTypeID: clone(details.CompanyResourceTypeID),
		description:           details.Description,
		location:              details.Location,
		maxReservationMinutes: clone(details.MaxReservationMinutes),
		availableFrom:         clone(details.AvailableFrom),
		availableTo:           clone(details.AvailableTo),
		status:                StatusActive,
		createdAt:             now,
		updatedAt:             now,
	}, nil
}

func (r *Resource) ID() uuid.UUID                     { return r.id }
func (r *Resource) CompanyID() uuid.UUID              { return r.companyID }
func (r *Resource) Name() string                      { return r.name }
func (r *Resource) Type() ResourceType                { return r.resourceType }
func (r *Resource) CompanyResourceTypeID() *uuid.UUID { return clone(r.companyResourceTypeID) }
func (r *Resource) Description() string               { return r.description }
func (r *Resource) Location() string                  { return r.location }
func (r *Resource) MaxReservationMinutes() *int       { return clone(r.maxReservationMinutes) }
func (r *Resource) AvailableFrom() *string            { return clone(r.availableFrom) }
func (r *Resource) AvailableTo() *string              { return clone(r.availableTo) }
func (r *Resource) Status() ResourceStatus            { return r.status }
func (r *Resource) UnavailableFrom() *time.Time       { return clone(r.unavailableFrom) }
func (r *Resource) UnavailableTo() *time.Time         { return clone(r.unavailableTo) }
func (r *Resource) UnavailableReason() *string        { return clone(r.unavailableReason) }
func (r *Resource) CreatedAt() time.Time              { return r.createdAt }
func (r *Resource) UpdatedAt() time.Time              { return r.updatedAt }

type Props struct {
	ID                    uuid.UUID
	CompanyID             uuid.UUID
	Name                  string
	ResourceType          string
	CompanyResourceTypeID *uuid.UUID
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
		companyResourceTypeID: clone(p.CompanyResourceTypeID),
		description:           p.Description,
		location:              p.Location,
		maxReservationMinutes: clone(p.MaxReservationMinutes),
		availableFrom:         clone(p.AvailableFrom),
		availableTo:           clone(p.AvailableTo),
		status:                ResourceStatus(p.Status),
		unavailableFrom:       clone(p.UnavailableFrom),
		unavailableTo:         clone(p.UnavailableTo),
		unavailableReason:     clone(p.UnavailableReason),
		createdAt:             p.CreatedAt,
		updatedAt:             p.UpdatedAt,
	}
}
