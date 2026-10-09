package reservationport

import (
	"context"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

type GetMyRequest struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Page      int
	Limit     int
}

type GetMyResponse struct {
	Reservations []reservation.Reservation
	Total        int64
}

type GetMy interface {
	Execute(context.Context, *GetMyRequest) (*GetMyResponse, error)
}
