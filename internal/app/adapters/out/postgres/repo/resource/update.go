package resourcerepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"
)

func (r *ResourceRepo) Update(ctx context.Context, res *resource.Resource) error {
	q := `UPDATE resources SET name = $1, type = $2, description = $3, location = $4, max_reservation_minutes = $5, available_from = $6, available_to = $7, status = $8, unavailable_from = $9, unavailable_to = $10, unavailable_reason = $11, company_resource_type_id = $12, updated_at = $13 WHERE id = $14`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		res.Name(),
		res.Type().String(),
		res.Description(),
		res.Location(),
		res.MaxReservationMinutes(),
		res.AvailableFrom(),
		res.AvailableTo(),
		res.Status().String(),
		res.UnavailableFrom(),
		res.UnavailableTo(),
		res.UnavailableReason(),
		res.CompanyResourceTypeID(),
		res.UpdatedAt(),
		res.ID(),
	)
	if err != nil {
		return fmt.Errorf("exec update resource: %w", err)
	}

	return nil
}
