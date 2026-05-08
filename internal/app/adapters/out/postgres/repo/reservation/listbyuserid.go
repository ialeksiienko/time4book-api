package reservationrepo

import (
	"context"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*booking.Booking, int64, error) {
	filter := booking.ListFilter{
		UserID: &userID,
		Page:   page,
		Limit:  limit,
	}
	return r.List(ctx, filter)
}
