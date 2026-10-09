package authport

import (
	"context"
	outports "time4book/internal/app/core/ports/out"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Firstname       string `validate:"required"`
	Lastname        string `validate:"required"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	CompanyName     string `validate:"required"`
	CompanyNIP      *string
	CompanyAddress  *string
	CompanyIndustry *string
}

type RegisterResponse struct {
	UserID       uuid.UUID
	CompanyID    uuid.UUID
	AccessToken  *outports.Token
	RefreshToken *outports.Token
}

type Register interface {
	Execute(context.Context, *RegisterRequest) (*RegisterResponse, error)
}
