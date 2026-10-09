package companyresourcetypequeries

import (
	"context"
	"fmt"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/google/uuid"
)

type Get struct {
	repo     companyresourcetypeout.Repository
	userRepo userout.Repository
}

func (c *Get) Execute(ctx context.Context, req *port.GetRequest) (*port.GetResponse, error) {
	if req.CompanyID == uuid.Nil {
		return nil, user.ErrUnauthorized
	}
	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, err
	}
	if !initiator.IsActive() {
		return nil, user.ErrUnauthorized
	}
	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() && !initiator.Role().IsEmployee() {
			return nil, user.ErrUnauthorized
		}
	}
	items, err := c.repo.Get(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company resource types: %w", err)
	}
	return &port.GetResponse{Items: items}, nil
}
