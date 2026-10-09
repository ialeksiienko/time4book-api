package authport

import (
	"context"
	domain "time4book/internal/app/core/domain/model/auth"

	"github.com/google/uuid"
)

type Repository interface {
	InsertUserCredentials(ctx context.Context, c *domain.Credentials) error
	RetrieveCredentialsByEmail(ctx context.Context, email string) (*domain.Credentials, error)
	InsertUserSession(ctx context.Context, s *domain.Session) error
	RetrieveSessionByRefreshToken(ctx context.Context, token string) (*domain.Session, error)
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteSessionByRefreshToken(ctx context.Context, token string) error
}
