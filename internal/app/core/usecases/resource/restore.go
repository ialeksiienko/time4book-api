package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type Restore struct {
	userRepo     userout.Repository
	resourceRepo resourceout.Repository
	log          *slog.Logger
}

func newRestore(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	l *slog.Logger,
) *Restore {
	return &Restore{
		userRepo:     urepo,
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *Restore) Execute(ctx context.Context, req *port.RestoreRequest) (*port.RestoreResponse, error) {
	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	res, err := c.resourceRepo.Retrieve(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if req.CompanyID != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	updatedRes := res.Restore()

	if err := c.resourceRepo.Update(ctx, updatedRes); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &port.RestoreResponse{}, nil
}
