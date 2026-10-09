package reservationrepo

import (
	"context"
	"errors"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"

	"github.com/jackc/pgconn"
)

func (r *ReservationRepo) Create(ctx context.Context, b *reservation.Reservation) error {
	q := `INSERT INTO reservations (id, user_id, company_id, resource_id, start_date, end_date, description, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		b.ID(),
		b.UserID(),
		b.CompanyID(),
		b.ResourceID(),
		b.StartDate(),
		b.EndDate(),
		b.Description(),
		b.Status().String(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" && pgErr.ConstraintName == "reservations_active_no_overlap" {
			return reservationout.ErrTimeSlotBooked
		}
		return fmt.Errorf("exec insert reservation: %w", err)
	}
	return nil
}
