package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepo interface {
	InsertUserCredentials(ctx context.Context, c *Credentials) error
	GetCredentialsByEmail(ctx context.Context, email string) (*Credentials, error)
	//
	InsertUserSession(ctx context.Context, s *Session) error
	GetSessionByRefreshToken(ctx context.Context, token string) (*Session, error)
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteSessionByRefreshToken(ctx context.Context, token string) error
}

