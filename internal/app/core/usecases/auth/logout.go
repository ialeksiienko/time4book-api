package authcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"

	"github.com/google/uuid"
)

type LogoutRequest struct {
	UserID uuid.UUID
}

type LogoutResponse struct{}

type Logout struct {
	authRepo auth.AuthRepo
	log      *slog.Logger
}

func newLogout(
	arepo auth.AuthRepo,
	l *slog.Logger,
) *Logout {
	return &Logout{
		authRepo: arepo,
		log:      l,
	}
}

func (c *Logout) Execute(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	err := c.authRepo.DeleteSessionsByUserID(ctx, req.UserID)
	if err != nil {
		c.log.Error("delete sessions by user id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("delete sessions: %w", err)
	}

	return &LogoutResponse{}, nil
}
