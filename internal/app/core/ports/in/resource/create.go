package resourceport

import (
	"context"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID           uuid.UUID
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

type CreateResponse struct {
	ResourceID uuid.UUID
}

type Create interface {
	Execute(context.Context, *CreateRequest) (*CreateResponse, error)
}
