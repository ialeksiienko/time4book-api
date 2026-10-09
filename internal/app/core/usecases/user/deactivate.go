package usercommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/user"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type Deactivate struct {
	userRepo userout.Repository
	log      *slog.Logger
}

func newDeactivate(
	urepo userout.Repository,
	l *slog.Logger,
) *Deactivate {
	return &Deactivate{
		userRepo: urepo,
		log:      l,
	}
}

func (c *Deactivate) Execute(ctx context.Context, req *port.DeactivateRequest) (*port.DeactivateResponse, error) {
	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	target, err := c.userRepo.Retrieve(ctx, req.TargetID)
	if err != nil {
		return nil, fmt.Errorf("get target: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if req.CompanyID == uuid.Nil || initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if target.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	target.Deactivate()

	if err := c.userRepo.Update(ctx, target); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &port.DeactivateResponse{}, nil
}
