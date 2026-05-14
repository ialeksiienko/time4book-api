package companyresourcetypecommands

import (
	"context"
	"fmt"
	"log/slog"

	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ID          uuid.UUID
	Name        string `validate:"required"`
	IconKey     string `validate:"required"`
}

type UpdateResponse struct{}

type Update struct {
	repo  companyresourcetype.Repo
	urepo user.UserRepo
	v     *validator.Facade
	log   *slog.Logger
}

func newUpdate(
	repo companyresourcetype.Repo,
	urepo user.UserRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Update {
	return &Update{repo: repo, urepo: urepo, v: v, log: l}
}

func (c *Update) Execute(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	if err := c.v.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.urepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	t, err := c.repo.ByIDAndCompany(ctx, req.ID, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company resource type: %w", err)
	}

	if err := t.Update(req.Name, req.IconKey); err != nil {
		return nil, err
	}

	if err := c.repo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update company resource type: %w", err)
	}

	return &UpdateResponse{}, nil
}
