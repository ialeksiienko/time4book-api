package authrepo

import (
	"context"
	"fmt"
	"time"
	"time4book/internal/app/core/domain/model/auth"

	"github.com/google/uuid"
)

func (r *AuthRepo) GetSessionByRefreshToken(ctx context.Context, token string) (*auth.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, created_at, updated_at
		FROM user_sessions
		WHERE refresh_token = $1
	`

	var (
		id           uuid.UUID
		userID       uuid.UUID
		refreshToken string
		expiresAt    time.Time
		createdAt    time.Time
		updatedAt    *time.Time
	)

	err := r.db.QueryRow(ctx, query, token).Scan(
		&id,
		&userID,
		&refreshToken,
		&expiresAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user session: %w", err)
	}

	return auth.RestoreSession(
		id,
		userID,
		refreshToken,
		expiresAt,
		createdAt,
		updatedAt,
	), nil
}
