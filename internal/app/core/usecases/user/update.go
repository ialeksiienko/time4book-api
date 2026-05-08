package usercommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID uuid.UUID
	TargetID    uuid.UUID
	Firstname   *string
	Lastname    *string
	Role        *string
	Status      *string
}

type UpdateResponse struct{}

type Update struct {
	userRepo  user.UserRepo
	validator *validator.Facade
	log       *slog.Logger
}

func newUpdate(
	urepo user.UserRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Update {
	return &Update{
		userRepo:  urepo,
		validator: v,
		log:       l,
	}
}

func (c *Update) Execute(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	target, err := c.userRepo.ByID(ctx, req.TargetID)
	if err != nil {
		return nil, fmt.Errorf("get target: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || target.CompanyID() == nil || *initiator.CompanyID() != *target.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() && initiator.ID() != target.ID() {
			return nil, user.ErrUnauthorized
		}
	}

	if req.Firstname != nil || req.Lastname != nil {
		target.UpdateProfile(req.Firstname, req.Lastname)
	}

	if req.Role != nil {
		newRole, err := user.NewRole(user.RoleKeyFromString(*req.Role), *req.Role)
		if err != nil {
			return nil, fmt.Errorf("invalid role: %w", err)
		}
		if !initiator.CanCreateUserWithRole(newRole) {
			return nil, user.ErrUnauthorized
		}
		target.ChangeRole(newRole)
	}

	if req.Status != nil {
		target.ChangeStatus(user.UserStatus(*req.Status))
	}

	if err := c.userRepo.Update(ctx, target); err != nil {
		return nil, fmt.Errorf("update user repo: %w", err)
	}

	return &UpdateResponse{}, nil
}

