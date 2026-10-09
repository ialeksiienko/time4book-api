package authport

import (
	"context"
	outports "time4book/internal/app/core/ports/out"
)

type RefreshRequest struct {
	RefreshToken string `validate:"required"`
}

type RefreshResponse struct {
	AccessToken  *outports.Token
	RefreshToken *outports.Token
}

type Refresh interface {
	Execute(context.Context, *RefreshRequest) (*RefreshResponse, error)
}
