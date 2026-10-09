package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type Delete struct {
	repo     companyresourcetypeout.Repository
	userRepo userout.Repository
}

func (c *Delete) Execute(ctx context.Context, req *port.DeleteRequest) (*port.DeleteResponse, error) {
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
	return &port.DeleteResponse{}, nil
}
