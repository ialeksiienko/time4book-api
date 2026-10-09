package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID, CompanyID, ID uuid.UUID
	Name, IconKey              string
}
type UpdateResponse struct{}
type Update struct {
	repo     companyresourcetype.Repo
	userRepo user.UserRepo
}

func (c *Update) Execute(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	typ, err := c.repo.Retrieve(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("retrieve company resource type: %w", err)
	}
	if req.CompanyID != uuid.Nil && req.CompanyID != typ.CompanyID() {
		return nil, user.ErrUnauthorized
	}
	if err := authorize(ctx, c.userRepo, req.InitiatorID, typ.CompanyID(), true); err != nil {
		return nil, err
	}
	if err := typ.Update(req.Name, req.IconKey); err != nil {
		return nil, err
	}
	if err := c.repo.Update(ctx, typ); err != nil {
		return nil, fmt.Errorf("update company resource type: %w", err)
	}
	return &UpdateResponse{}, nil
}
