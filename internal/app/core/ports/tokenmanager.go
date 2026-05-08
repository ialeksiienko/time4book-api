package ports

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type TokenManager interface {
	GenerateToken(userID uuid.UUID, role string, tokenType string) (*Token, error)
	ValidateToken(token string) (uuid.UUID, string, error)
}
