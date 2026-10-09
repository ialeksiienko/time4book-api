package authport

import (
	"context"
	outports "time4book/internal/app/core/ports/out"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	UserID       uuid.UUID
	AccessToken  *outports.Token
	RefreshToken *outports.Token
}

type Login interface {
	Execute(context.Context, *LoginRequest) (*LoginResponse, error)
}
