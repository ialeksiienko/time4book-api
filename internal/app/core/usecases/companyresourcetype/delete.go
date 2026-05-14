package companyresourcetypecommands

import (
	"context"
	"fmt"
	"log/slog"

	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type DeleteRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	ID          uuid.UUID
}

type DeleteResponse struct{}

type Delete struct {
	repo  companyresourcetype.Repo
	urepo user.UserRepo
	log   *slog.Logger
}

func newDelete(
	repo companyresourcetype.Repo,
	urepo user.UserRepo,
	l *slog.Logger,
) *Delete {
	return &Delete{repo: repo, urepo: urepo, log: l}
}

func (c *Delete) Execute(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
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

	if _, err := c.repo.ByIDAndCompany(ctx, req.ID, req.CompanyID); err != nil {
		return nil, fmt.Errorf("get company resource type: %w", err)
	}

	usedBy, err := c.repo.CountResourcesUsing(ctx, req.ID, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("count resources using company resource type: %w", err)
	}
	if usedBy > 0 {
		return nil, ErrTypeInUse
	}

	if err := c.repo.Delete(ctx, req.ID, req.CompanyID); err != nil {
		return nil, fmt.Errorf("delete company resource type: %w", err)
	}

	return &DeleteResponse{}, nil
}
