package authrepo

import (
	"context"
	"fmt"
)

func (r *AuthRepo) DeleteSessionByRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_sessions WHERE refresh_token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete user session: %w", err)
	}

	return nil
}
