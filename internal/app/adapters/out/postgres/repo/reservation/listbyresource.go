package reservationrepo

import (
	"context"
	"time"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ListByResourceID(ctx context.Context, resourceID uuid.UUID, companyID uuid.UUID, from, to *time.Time, page, limit int) ([]*reservation.Reservation, int64, error) {
	filter := reservation.ListFilter{
		ResourceID: &resourceID,
		CompanyID:  &companyID,
		From:       from,
		To:         to,
		Page:       page,
		Limit:      limit,
	}
	return r.List(ctx, filter)
}
