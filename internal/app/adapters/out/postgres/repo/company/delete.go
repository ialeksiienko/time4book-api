package companyrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"

	"github.com/google/uuid"
)

func (r *CompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM companies WHERE id = $1`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("exec delete company: %w", err)
	}
	return nil
}
