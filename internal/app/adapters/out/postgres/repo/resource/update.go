package resourcerepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/resource"
)

func (r *ResourceRepo) Update(ctx context.Context, res *resource.Resource) error {
	q := `UPDATE resources SET name = $1, type = $2, company_resource_type_id = $3, description = $4, location = $5, max_reservation_minutes = $6, available_from = $7, available_to = $8, status = $9, unavailable_from = $10, unavailable_to = $11, unavailable_reason = $12, updated_at = $13 WHERE id = $14 AND company_id = $15`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q,
		res.Name(),
		res.Type().String(),
		res.CompanyResourceTypeID(),
		res.Description(),
		res.Location(),
		res.MaxReservationMinutes(),
		res.AvailableFrom(),
		res.AvailableTo(),
		res.Status().String(),
		res.UnavailableFrom(),
		res.UnavailableTo(),
		res.UnavailableReason(),
		res.UpdatedAt(),
		res.ID(),
		res.CompanyID(),
	)
	if err != nil {
		return fmt.Errorf("exec update resource: %w", err)
	}

	return nil
}
