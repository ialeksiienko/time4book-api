package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/resource"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type Update struct {
	userRepo        userout.Repository
	resourceRepo    resourceout.Repository
	companyTypeRepo companyresourcetypeout.Repository
	validator       *validator.Facade
	log             *slog.Logger
}

func newUpdate(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	companyTypeRepo companyresourcetypeout.Repository,
	v *validator.Facade,
	l *slog.Logger,
) *Update {
	return &Update{
		userRepo:        urepo,
		resourceRepo:    resrepo,
		companyTypeRepo: companyTypeRepo,
		validator:       v,
		log:             l,
	}
}

func (c *Update) Execute(ctx context.Context, req *port.UpdateRequest) (*port.UpdateResponse, error) {
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
	if req.CompanyID != uuid.Nil && req.CompanyID != res.CompanyID() {
		return nil, user.ErrUnauthorized
	}
	if err := resource.ValidateType(resource.ResourceType(req.Type), req.CompanyResourceTypeID); err != nil {
		return nil, err
	}
	if req.Type == resource.TypeCustom.String() {
		typ, err := c.companyTypeRepo.Retrieve(ctx, *req.CompanyResourceTypeID)
		if err != nil {
			return nil, fmt.Errorf("retrieve company resource type: %w", err)
		}
		if typ.CompanyID() != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
	}

	updatedRes, err := res.UpdateDetails(resource.Details{
		Name:                  req.Name,
		Type:                  resource.ResourceType(req.Type),
		CompanyResourceTypeID: req.CompanyResourceTypeID,
		Description:           req.Description,
		Location:              req.Location,
		MaxReservationMinutes: req.MaxReservationMinutes,
		AvailableFrom:         req.AvailableFrom,
		AvailableTo:           req.AvailableTo,
	})
	if err != nil {
		return nil, fmt.Errorf("update resource details: %w", err)
	}

	if err := c.resourceRepo.Update(ctx, updatedRes); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &port.UpdateResponse{}, nil
}
