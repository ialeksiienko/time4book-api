package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type Service struct {
	userRepo     userout.Repository
	resourceRepo resourceout.Repository
	validator    *validator.Facade
	log          *slog.Logger
}

func newService(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	v *validator.Facade,
	l *slog.Logger,
) *Service {
	return &Service{
		userRepo:     urepo,
		resourceRepo: resrepo,
		validator:    v,
		log:          l,
	}
}

func (c *Service) Execute(ctx context.Context, req *port.ServiceRequest) (*port.ServiceResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	res, err := c.resourceRepo.Retrieve(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if req.CompanyID != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	updatedRes, err := res.MarkInService(req.Reason, req.From, req.To)
	if err != nil {
		return nil, fmt.Errorf("mark resource in service: %w", err)
	}

	if err := c.resourceRepo.Update(ctx, updatedRes); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &port.ServiceResponse{}, nil
}
