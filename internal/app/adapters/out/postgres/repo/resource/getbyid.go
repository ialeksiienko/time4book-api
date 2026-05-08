package resourcerepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"

	"github.com/google/uuid"
)

func (r *ResourceRepo) ByID(ctx context.Context, id uuid.UUID) (*resource.Resource, error) {
	q := `SELECT id, company_id, name, type, description, location, max_reservation_minutes, available_from, available_to, status, unavailable_from, unavailable_to, unavailable_reason, created_at, updated_at 
          FROM resources WHERE id = $1`

	var row struct {
		ID                    uuid.UUID
		CompanyID             uuid.UUID
		Name                  string
		Type                  string
		Description           string
		Location              string
		MaxReservationMinutes *int
		AvailableFrom         *string
		AvailableTo           *string
		Status                string
		UnavailableFrom       *time.Time
		UnavailableTo         *time.Time
		UnavailableReason     *string
		CreatedAt             time.Time
		UpdatedAt             time.Time
	}

	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, id).Scan(
		&row.ID,
		&row.CompanyID,
		&row.Name,
		&row.Type,
		&row.Description,
		&row.Location,
		&row.MaxReservationMinutes,
		&row.AvailableFrom,
		&row.AvailableTo,
		&row.Status,
		&row.UnavailableFrom,
		&row.UnavailableTo,
		&row.UnavailableReason,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan resource: %w", err)
	}

	return resource.Reconstitute(&resource.Props{
		ID:                    row.ID,
		CompanyID:             row.CompanyID,
		Name:                  row.Name,
		ResourceType:          row.Type,
		Description:           row.Description,
		Location:              row.Location,
		MaxReservationMinutes: row.MaxReservationMinutes,
		AvailableFrom:         row.AvailableFrom,
		AvailableTo:           row.AvailableTo,
		Status:                row.Status,
		UnavailableFrom:       row.UnavailableFrom,
		UnavailableTo:         row.UnavailableTo,
		UnavailableReason:     row.UnavailableReason,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}), nil
}
