package resource

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	id                    uuid.UUID
	companyID             uuid.UUID
	name                  string
	resourceType          ResourceType
	companyResourceTypeID *uuid.UUID
	customTypeName        *string
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

var ErrInvalidType = errors.New("invalid resource type")

func ValidateType(resType ResourceType, companyResourceTypeID *uuid.UUID) error {
	if resType == TypeCustom && companyResourceTypeID == nil {
		return fmt.Errorf("%w: custom resource type id is required", ErrInvalidType)
	}
	if resType != TypeCustom && companyResourceTypeID != nil {
		return fmt.Errorf("%w: custom resource type id is only for custom resources", ErrInvalidType)
	}
	switch resType {
	case TypeRoom, TypeCar, TypeMultimedia, TypeEquipment, TypeCustom:
		return nil
	default:
		return ErrInvalidType
	}
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
	if err := ValidateType(resType, companyResourceTypeID); err != nil {
		return nil, err
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

func (r *Resource) IsBookable() bool { return r.status == StatusActive }

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
