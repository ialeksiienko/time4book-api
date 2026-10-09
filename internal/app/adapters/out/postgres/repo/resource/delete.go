package resourcerepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"

	"github.com/google/uuid"
)

func (r *ResourceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM resources WHERE id = $1`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("execute delete resource: %w", err)
	}

	return nil
}
