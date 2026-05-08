package authrepo

import (
	"context"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"
	"time4book/internal/app/core/domain/model/auth"
)

func (r *AuthRepo) InsertUserSession(ctx context.Context, s *auth.Session) error {
	q := `INSERT INTO user_sessions (id, user_id, refresh_token, expires_at)
          VALUES ($1, $2, $3, $4)`

	_, err := postgres.ExtractQuerier(ctx, r.db).Exec(ctx, q, s.ID(), s.UserID(), s.RefreshToken(), s.ExpiresAt())
	if err != nil {
		return fmt.Errorf("exec insert user credentials: %w", err)
	}
	return nil
}
