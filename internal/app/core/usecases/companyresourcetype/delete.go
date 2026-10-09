package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type DeleteRequest struct{ InitiatorID, CompanyID, ID uuid.UUID }
type DeleteResponse struct{}
type Delete struct {
	repo     companyresourcetype.Repo
	userRepo user.UserRepo
}

func (c *Delete) Execute(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
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
	if err := c.repo.Delete(ctx, req.ID, typ.CompanyID()); err != nil {
		return nil, fmt.Errorf("delete company resource type: %w", err)
	}
	return &DeleteResponse{}, nil
}
