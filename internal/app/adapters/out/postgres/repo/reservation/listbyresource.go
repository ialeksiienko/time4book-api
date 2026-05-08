package reservationrepo

import (
	"context"
	"time"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ListByResourceID(ctx context.Context, resourceID uuid.UUID, from, to *time.Time, page, limit int) ([]*booking.Booking, int64, error) {
	filter := booking.ListFilter{
		ResourceID: &resourceID,
		From:       from,
		To:         to,
		Page:       page,
		Limit:      limit,
	}
	return r.List(ctx, filter)
}
