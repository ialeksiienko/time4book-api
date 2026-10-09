package reservationport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID uuid.UUID `validate:"required"`
	CompanyID   uuid.UUID `validate:"required"`
	ResourceID  uuid.UUID `validate:"required"`
	StartDate   time.Time `validate:"required"`
	EndDate     time.Time `validate:"required"`
	Description *string
}

type CreateResponse struct {
	ReservationID uuid.UUID
}

type Create interface {
	Execute(context.Context, *CreateRequest) (*CreateResponse, error)
}
