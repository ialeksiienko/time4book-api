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

type Delete struct {
	userRepo     userout.Repository
	resourceRepo resourceout.Repository
	log          *slog.Logger
}

func newDelete(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	l *slog.Logger,
) *Delete {
	return &Delete{
		userRepo:     urepo,
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *Delete) Execute(ctx context.Context, req *port.DeleteRequest) (*port.DeleteResponse, error) {
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

	if err := c.resourceRepo.Delete(ctx, req.ResourceID); err != nil {
		return nil, fmt.Errorf("delete resource: %w", err)
	}

	return &port.DeleteResponse{}, nil
}
