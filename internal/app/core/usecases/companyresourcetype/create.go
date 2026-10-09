package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID, CompanyID uuid.UUID
	Name, IconKey          string
}
type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}
type Create struct {
	repo     companyresourcetype.Repo
	userRepo user.UserRepo
}

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if req.CompanyID == uuid.Nil {
		return nil, user.ErrUnauthorized
	}
	if err := authorize(ctx, c.userRepo, req.InitiatorID, req.CompanyID, true); err != nil {
		return nil, err
	}
	typ, err := companyresourcetype.NewCompanyResourceType(req.CompanyID, req.Name, req.IconKey)
	if err != nil {
		return nil, err
	}
	if err := c.repo.Create(ctx, typ); err != nil {
		return nil, fmt.Errorf("create company resource type: %w", err)
	}
	return &CreateResponse{ID: typ.ID()}, nil
}
