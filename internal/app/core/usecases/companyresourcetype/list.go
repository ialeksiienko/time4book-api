package companyresourcetypecommands

import (
	"context"
	"fmt"
	"log/slog"

	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type ListRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
}

type ListResponse struct {
	Items []*companyresourcetype.CompanyResourceType
}

type List struct {
	repo  companyresourcetype.Repo
	urepo user.UserRepo
	log   *slog.Logger
}

func newList(repo companyresourcetype.Repo, urepo user.UserRepo, l *slog.Logger) *List {
	return &List{repo: repo, urepo: urepo, log: l}
}

func (q *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	initiator, err := q.urepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() && !initiator.Role().IsEmployee() {
			return nil, user.ErrUnauthorized
		}
	}

	items, err := q.repo.ListByCompany(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("list company resource types: %w", err)
	}

	return &ListResponse{Items: items}, nil
}
