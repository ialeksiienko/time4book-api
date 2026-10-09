package reservationport

import (
	"context"
	"time"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

type GetRequest struct {
	CompanyID  *uuid.UUID
	UserID     *uuid.UUID
	ResourceID *uuid.UUID
	Status     *string
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type GetResponse struct {
	Reservations []reservation.Reservation
	Total        int64
}

type Get interface {
	Execute(context.Context, *GetRequest) (*GetResponse, error)
}
