package usercommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type DeactivateRequest struct {
	InitiatorID uuid.UUID
	TargetID    uuid.UUID
}

type DeactivateResponse struct{}

type Deactivate struct {
	userRepo user.UserRepo
	log      *slog.Logger
}

func newDeactivate(
	urepo user.UserRepo,
	l *slog.Logger,
) *Deactivate {
	return &Deactivate{
		userRepo: urepo,
		log:      l,
	}
}

func (c *Deactivate) Execute(ctx context.Context, req *DeactivateRequest) (*DeactivateResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	target, err := c.userRepo.ByID(ctx, req.TargetID)
	if err != nil {
		return nil, fmt.Errorf("get target: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || target.CompanyID() == uuid.Nil || initiator.CompanyID() != target.CompanyID() {
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

	return &DeactivateResponse{}, nil
}
