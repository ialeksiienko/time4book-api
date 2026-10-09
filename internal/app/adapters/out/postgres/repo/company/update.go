package companyrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/company"
)

func (r *CompanyRepo) Update(ctx context.Context, c *company.Company) error {
	q := `UPDATE companies SET name = $1, nip = $2, address = $3, industry = $4, status = $5, updated_at = $6 WHERE id = $7`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, c.Name(), c.NIP(), c.Address(), c.Industry(), c.Status().String(), c.UpdatedAt(), c.ID())
	if err != nil {
		return fmt.Errorf("exec update company: %w", err)
	}
	return nil
}
