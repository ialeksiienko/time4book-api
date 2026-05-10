package reservationrepo

import (
	"context"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ListByUserID(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, page, limit int) ([]*reservation.Reservation, int64, error) {
	filter := reservation.ListFilter{
		UserID:    &userID,
		CompanyID: &companyID,
		Page:      page,
		Limit:     limit,
	}
	return r.List(ctx, filter)
}
