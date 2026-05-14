package authrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/auth"
)

func (r *AuthRepo) InsertUserCredentials(ctx context.Context, c *auth.Credentials) error {
	q := `INSERT INTO user_credentials (user_id, email, password_hash)
        VALUES ($1, $2, $3)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, c.UserID(), c.Email(), c.PasswordHash())
	if err != nil {
		return fmt.Errorf("exec insert user credentials: %w", err)
	}
	return nil
}

func (r *AuthRepo) UpsertUserCredentials(ctx context.Context, c *auth.Credentials) error {
	q := `INSERT INTO user_credentials (user_id, email, password_hash)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id)
        DO UPDATE SET
            email = EXCLUDED.email,
            password_hash = EXCLUDED.password_hash`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, c.UserID(), c.Email(), c.PasswordHash())
	if err != nil {
		return fmt.Errorf("exec upsert user credentials: %w", err)
	}
	return nil
}
