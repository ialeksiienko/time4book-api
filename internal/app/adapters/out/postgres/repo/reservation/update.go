package reservationrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/booking"
)

func (r *ReservationRepo) Update(ctx context.Context, b *booking.Booking) error {
	q := `UPDATE reservations SET start_date = $1, end_date = $2, description = $3, status = $4, updated_at = $5 WHERE id = $6`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		b.StartDate(),
		b.EndDate(),
		b.Description(),
		b.Status().String(),
		b.UpdatedAt(),
		b.ID(),
	)
	if err != nil {
		return fmt.Errorf("exec update reservation: %w", err)
	}

	return nil
}
