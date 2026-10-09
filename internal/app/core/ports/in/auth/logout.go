package authport

import (
	"context"

	"github.com/google/uuid"
)

type LogoutRequest struct {
	UserID uuid.UUID
}

type LogoutResponse struct{}

type Logout interface {
	Execute(context.Context, *LogoutRequest) (*LogoutResponse, error)
}
