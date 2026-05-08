package authrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"

	"github.com/google/uuid"
)

func (r *AuthRepo) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	q := `DELETE FROM user_sessions WHERE user_id = $1`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("execute delete session: %w", err)
	}

	return nil
}
