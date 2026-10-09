package resource

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Details struct {
	Name                  string
	Type                  ResourceType
	CompanyResourceTypeID *uuid.UUID
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
}

func ValidateType(resType ResourceType, companyResourceTypeID *uuid.UUID) error {
	if resType == TypeCustom && (companyResourceTypeID == nil || *companyResourceTypeID == uuid.Nil) {
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

func validateDetails(details Details) error {
	if strings.TrimSpace(details.Name) == "" {
		return fmt.Errorf("%w: resource name cannot be empty", ErrInvalidDetails)
	}
	if err := ValidateType(details.Type, details.CompanyResourceTypeID); err != nil {
		return err
	}
	if details.MaxReservationMinutes != nil && *details.MaxReservationMinutes <= 0 {
		return fmt.Errorf("%w: maximum reservation minutes must be positive", ErrInvalidDetails)
	}
	if (details.AvailableFrom == nil) != (details.AvailableTo == nil) {
		return fmt.Errorf("%w: both daily hours are required", ErrInvalidDetails)
	}
	if details.AvailableFrom != nil {
		from, err := time.Parse("15:04", *details.AvailableFrom)
		if err != nil || from.Format("15:04") != *details.AvailableFrom {
			return fmt.Errorf("%w: available from must use HH:mm", ErrInvalidDetails)
		}
		to, err := time.Parse("15:04", *details.AvailableTo)
		if err != nil || to.Format("15:04") != *details.AvailableTo {
			return fmt.Errorf("%w: available to must use HH:mm", ErrInvalidDetails)
		}
		if !to.After(from) {
			return fmt.Errorf("%w: available to must be after available from", ErrInvalidDetails)
		}
	}
	return nil
}

func (r *Resource) UpdateDetails(details Details) (*Resource, error) {
	if err := validateDetails(details); err != nil {
		return nil, err
	}
	updated := *r
	updated.name = details.Name
	updated.resourceType = details.Type
	updated.companyResourceTypeID = clone(details.CompanyResourceTypeID)
	updated.description = details.Description
	updated.location = details.Location
	updated.maxReservationMinutes = clone(details.MaxReservationMinutes)
	updated.availableFrom = clone(details.AvailableFrom)
	updated.availableTo = clone(details.AvailableTo)
	updated.updatedAt = time.Now().UTC()
	return &updated, nil
}
