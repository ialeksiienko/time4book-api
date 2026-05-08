package authrepo

import (
	"context"
	"errors"
	"fmt"
	"time4book/internal/app/adapters/out/postgres"

	"time4book/internal/app/core/domain/model/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (r *AuthRepo) GetCredentialsByEmail(ctx context.Context, email string) (*auth.Credentials, error) {
	q := `SELECT user_id, email, password_hash FROM user_credentials WHERE email = $1`

	var credsRow struct {
		UserID       uuid.UUID
		Email        string
		PasswordHash string
	}

	err := postgres.ExtractQuerier(ctx, r.db).QueryRow(ctx, q, email).Scan(&credsRow.UserID, &credsRow.Email, &credsRow.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user credentials not found: %w", err)
		}
		return nil, fmt.Errorf("scan user credentials: %w", err)
	}

	return auth.ReconstituteCredentials(&auth.CredentialsProps{
		UserID:       credsRow.UserID,
		Email:        credsRow.Email,
		PasswordHash: credsRow.PasswordHash,
	}), nil
}
