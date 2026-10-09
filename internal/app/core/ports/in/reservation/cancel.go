package reservationport

import (
	"context"

	"github.com/google/uuid"
)

type CancelRequest struct {
	InitiatorID   uuid.UUID
	CompanyID     uuid.UUID
	ReservationID uuid.UUID
}

type CancelResponse struct{}

type Cancel interface {
	Execute(context.Context, *CancelRequest) (*CancelResponse, error)
}
