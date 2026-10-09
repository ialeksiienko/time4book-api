package resourcerepo

import (
	"context"
	"errors"
	"fmt"
	"time"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (r *ResourceRepo) Retrieve(ctx context.Context, id uuid.UUID) (*resource.Resource, error) {
	q := `
SELECT r.id, r.company_id, r.name, r.type, r.description, r.location, r.max_reservation_minutes, r.available_from, r.available_to, r.status, r.unavailable_from, r.unavailable_to, r.unavailable_reason, r.created_at, r.updated_at,
       r.company_resource_type_id
FROM resources r
WHERE r.id = $1`

	var row struct {
		ID                      uuid.UUID
		CompanyID               uuid.UUID
		Name                    string
		Type                    string
		Description             string
		Location                string
		MaxReservationMinutes   *int
		AvailableFrom           *string
		AvailableTo             *string
		Status                  string
		UnavailableFrom         *time.Time
		UnavailableTo           *time.Time
		UnavailableReason       *string
		CreatedAt               time.Time
		UpdatedAt               time.Time
		CompanyResourceTypeIDNu uuid.NullUUID
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
		&row.CompanyResourceTypeIDNu,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, resourceout.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan resource: %w", err)
	}

	var companyCRTID *uuid.UUID
	if row.CompanyResourceTypeIDNu.Valid {
		u := row.CompanyResourceTypeIDNu.UUID
		companyCRTID = &u
	}
	return resource.Reconstitute(&resource.Props{
		ID:                    row.ID,
		CompanyID:             row.CompanyID,
		Name:                  row.Name,
		ResourceType:          row.Type,
		CompanyResourceTypeID: companyCRTID,
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
