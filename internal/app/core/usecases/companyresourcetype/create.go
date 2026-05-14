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

type CreateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	Name        string `validate:"required"`
	IconKey     string `validate:"required"`
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

type Create struct {
	repo  companyresourcetype.Repo
	urepo user.UserRepo
	v     *validator.Facade
	log   *slog.Logger
}

func newCreate(
	repo companyresourcetype.Repo,
	urepo user.UserRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{repo: repo, urepo: urepo, v: v, log: l}
}

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
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

	t, err := companyresourcetype.NewCompanyResourceType(req.CompanyID, req.Name, req.IconKey)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create company resource type: %w", err)
	}

	return &CreateResponse{ID: t.ID()}, nil
}
