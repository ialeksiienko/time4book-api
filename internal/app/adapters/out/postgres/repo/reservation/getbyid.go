package reservationrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

func (r *ReservationRepo) ByID(ctx context.Context, id uuid.UUID) (*reservation.Reservation, error) {
	q := `SELECT id, user_id, company_id, resource_id, start_date, end_date, description, status, created_at, updated_at 
          FROM reservations WHERE id = $1`

	var row struct {
		ID          uuid.UUID
		UserID      uuid.UUID
		CompanyID   uuid.UUID
		ResourceID  uuid.UUID
		StartDate   time.Time
		EndDate     time.Time
		Description *string
		Status      string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, id).Scan(
		&row.ID,
		&row.UserID,
		&row.CompanyID,
		&row.ResourceID,
		&row.StartDate,
		&row.EndDate,
		&row.Description,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan reservation: %w", err)
	}

	return reservation.Reconstitute(&reservation.Props{
		ID:          row.ID,
		UserID:      row.UserID,
		CompanyID:   row.CompanyID,
		ResourceID:  row.ResourceID,
		Description: row.Description,
		StartDate:   row.StartDate,
		EndDate:     row.EndDate,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}), nil
}
