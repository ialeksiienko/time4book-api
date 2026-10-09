package reservationport

import (
	"context"
	"time"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

type GetByResourceRequest struct {
	ResourceID uuid.UUID
	CompanyID  uuid.UUID
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type GetByResourceResponse struct {
	Reservations []reservation.Reservation
	Total        int64
}

type GetByResource interface {
	Execute(context.Context, *GetByResourceRequest) (*GetByResourceResponse, error)
}
