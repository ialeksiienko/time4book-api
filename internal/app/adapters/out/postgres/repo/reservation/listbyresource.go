package reservationrepo

import (
	"context"
	"time"
	"time4book/internal/app/core/domain/model/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) GetByResourceID(ctx context.Context, resourceID uuid.UUID, companyID uuid.UUID, from, to *time.Time, page, limit int) ([]*reservation.Reservation, int64, error) {
	filter := reservationout.Filter{
		ResourceID: &resourceID,
		CompanyID:  &companyID,
		From:       from,
		To:         to,
		Page:       page,
		Limit:      limit,
	}
	return r.Get(ctx, filter)
}
