package outports

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

type Tokens struct {
	AccessToken  *Token
	RefreshToken *Token
}

type TokenManager interface {
	GenerateTokens(userID uuid.UUID, role string) (*Tokens, error)
	ValidateToken(token string) (uuid.UUID, string, error)
}
