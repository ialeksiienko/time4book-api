package reservationrepo

import (
	"context"
	"time4book/internal/app/core/domain/model/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) GetByUserID(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, page, limit int) ([]*reservation.Reservation, int64, error) {
	filter := reservationout.Filter{
		UserID:    &userID,
		CompanyID: &companyID,
		Page:      page,
		Limit:     limit,
	}
	return r.Get(ctx, filter)
}
