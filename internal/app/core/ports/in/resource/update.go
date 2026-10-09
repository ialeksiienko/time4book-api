package resourceport

import (
	"context"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID           uuid.UUID
	ResourceID            uuid.UUID
	CompanyID             uuid.UUID
	Name                  string `validate:"required"`
	Type                  string `validate:"required"`
	CompanyResourceTypeID *uuid.UUID
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
}

type UpdateResponse struct{}

type Update interface {
	Execute(context.Context, *UpdateRequest) (*UpdateResponse, error)
}
