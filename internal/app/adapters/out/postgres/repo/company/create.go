package companyrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/company"
)

func (r *CompanyRepo) Create(ctx context.Context, c *company.Company) error {
	q := `INSERT INTO companies (id, owner_id, name, nip, industry, address, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, c.ID(), c.OwnerID(), c.Name(), c.NIP(), c.Industry(), c.Address(), c.Status().String(), c.CreatedAt())
	if err != nil {
		return fmt.Errorf("exec insert company: %w", err)
	}
	return nil
}
