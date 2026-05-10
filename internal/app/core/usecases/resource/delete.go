package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type DeleteRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ResourceID  uuid.UUID
}

type DeleteResponse struct{}

type Delete struct {
	userRepo     user.UserRepo
	resourceRepo resource.ResourceRepo
	log          *slog.Logger
}

func newDelete(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	l *slog.Logger,
) *Delete {
	return &Delete{
		userRepo:     urepo,
		resourceRepo: resrepo,
		log:          l,
	}
}

func (c *Delete) Execute(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
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

	return &DeleteResponse{}, nil
}
