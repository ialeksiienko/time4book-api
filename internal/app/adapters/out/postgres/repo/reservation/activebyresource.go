package reservationrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ActiveByResourceIDInRange(ctx context.Context, resourceID uuid.UUID, from, to time.Time, excludeID *uuid.UUID) ([]*booking.Booking, error) {
	q := `SELECT id, user_id, resource_id, start_date, end_date, description, status, created_at, updated_at 
          FROM reservations 
          WHERE resource_id = $1 
            AND status = 'active'
            AND start_date < $3 AND end_date > $2`

	args := []interface{}{resourceID, from, to}
	if excludeID != nil {
		q += ` AND id != $4`
		args = append(args, *excludeID)
	}

	rows, err := postgres.ExtractQuerier(ctx, r.db).Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query active reservations: %w", err)
	}
	defer rows.Close()

	var res []*booking.Booking
	for rows.Next() {
		var row struct {
			ID          uuid.UUID
			UserID      uuid.UUID
			ResourceID  uuid.UUID
			StartDate   time.Time
			EndDate     time.Time
			Description *string
			Status      string
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		if err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.ResourceID,
			&row.StartDate,
			&row.EndDate,
			&row.Description,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reservation: %w", err)
		}

		res = append(res, booking.Reconstitute(&booking.Props{
			ID:          row.ID,
			UserID:      row.UserID,
			ResourceID:  row.ResourceID,
			Description: row.Description,
			StartDate:   row.StartDate,
			EndDate:     row.EndDate,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}))
	}

	return res, nil
}
