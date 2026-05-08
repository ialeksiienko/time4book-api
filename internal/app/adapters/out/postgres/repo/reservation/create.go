package reservationrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/booking"
)

func (r *ReservationRepo) Create(ctx context.Context, b *booking.Booking) error {
	q := `INSERT INTO reservations (id, user_id, resource_id, start_date, end_date, description, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		b.ID(),
		b.UserID(),
		b.ResourceID(),
		b.StartDate(),
		b.EndDate(),
		b.Description(),
		b.Status().String(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("exec insert reservation: %w", err)
	}
	return nil
}
