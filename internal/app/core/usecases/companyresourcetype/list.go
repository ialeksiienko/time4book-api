package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type ListRequest struct{ InitiatorID, CompanyID uuid.UUID }
type ListResponse struct {
	Items []*companyresourcetype.CompanyResourceType
}
type List struct {
	repo     companyresourcetype.Repo
	userRepo user.UserRepo
}

func (c *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	if req.CompanyID == uuid.Nil {
		return nil, user.ErrUnauthorized
	}
	if err := authorize(ctx, c.userRepo, req.InitiatorID, req.CompanyID, false); err != nil {
		return nil, err
	}
	items, err := c.repo.Get(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company resource types: %w", err)
	}
	return &ListResponse{Items: items}, nil
}
