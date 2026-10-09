package companyresourcetypecommands

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type Create struct {
	repo     companyresourcetypeout.Repository
	userRepo userout.Repository
}

func (c *Create) Execute(ctx context.Context, req *port.CreateRequest) (*port.CreateResponse, error) {
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
	return &port.CreateResponse{ID: typ.ID()}, nil
}
