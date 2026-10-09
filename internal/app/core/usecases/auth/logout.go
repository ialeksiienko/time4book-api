package authcommands

import (
	"context"
	"fmt"
	"log/slog"
	port "time4book/internal/app/core/ports/in/auth"
	authout "time4book/internal/app/core/ports/out/auth"
)

type Logout struct {
	authRepo authout.Repository
	log      *slog.Logger
}

func newLogout(
	arepo authout.Repository,
	l *slog.Logger,
) *Logout {
	return &Logout{
		authRepo: arepo,
		log:      l,
	}
}

func (c *Logout) Execute(ctx context.Context, req *port.LogoutRequest) (*port.LogoutResponse, error) {
	err := c.authRepo.DeleteSessionsByUserID(ctx, req.UserID)
	if err != nil {
		c.log.Error("delete sessions by user id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("delete sessions: %w", err)
	}

	return &port.LogoutResponse{}, nil
}
