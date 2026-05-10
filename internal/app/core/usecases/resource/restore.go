package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type RestoreRequest struct {
	InitiatorID uuid.UUID
	ResourceID  uuid.UUID
}

type RestoreResponse struct{}

type Restore struct {
	userRepo     user.UserRepo
	resourceRepo resource.ResourceRepo
	log          *slog.Logger
}

func newRestore(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	l *slog.Logger,
) *Restore {
	return &Restore{
		userRepo:     urepo,
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *Restore) Execute(ctx context.Context, req *RestoreRequest) (*RestoreResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	res, err := c.resourceRepo.ByID(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	res.Restore()

	if err := c.resourceRepo.Update(ctx, res); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &RestoreResponse{}, nil
}
