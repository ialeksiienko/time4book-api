package resourcerepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"
)

func (r *ResourceRepo) Create(ctx context.Context, res *resource.Resource) error {
	q := `INSERT INTO resources (id, company_id, name, type, description, location, max_reservation_minutes, available_from, available_to, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		res.ID(),
		res.CompanyID(),
		res.Name(),
		res.Type().String(),
		res.Description(),
		res.Location(),
		res.MaxReservationMinutes(),
		res.AvailableFrom(),
		res.AvailableTo(),
		res.Status().String(),
		res.CreatedAt(),
		res.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("exec insert resource: %w", err)
	}
	return nil
}
