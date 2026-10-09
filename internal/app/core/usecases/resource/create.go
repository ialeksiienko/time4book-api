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

type Create struct {
	userRepo        userout.Repository
	resourceRepo    resourceout.Repository
	companyTypeRepo companyresourcetypeout.Repository
	validator       *validator.Facade
	log             *slog.Logger
}

func newCreate(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	companyTypeRepo companyresourcetypeout.Repository,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{
		userRepo:        urepo,
		resourceRepo:    resrepo,
		companyTypeRepo: companyTypeRepo,
		validator:       v,
		log:             l,
	}
}

func (c *Create) Execute(ctx context.Context, req *port.CreateRequest) (*port.CreateResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	resType := resource.ResourceType(req.Type)
	if req.CompanyID == uuid.Nil {
		return nil, user.ErrUnauthorized
	}
	if err := resource.ValidateType(resType, req.CompanyResourceTypeID); err != nil {
		return nil, err
	}
	if resType == resource.TypeCustom {
		typ, err := c.companyTypeRepo.Retrieve(ctx, *req.CompanyResourceTypeID)
		if err != nil {
			return nil, fmt.Errorf("retrieve company resource type: %w", err)
		}
		if typ.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
	}

	res, err := resource.NewResource(req.CompanyID, resource.Details{
		Name:                  req.Name,
		Type:                  resType,
		CompanyResourceTypeID: req.CompanyResourceTypeID,
		Description:           req.Description,
		Location:              req.Location,
		MaxReservationMinutes: req.MaxReservationMinutes,
		AvailableFrom:         req.AvailableFrom,
		AvailableTo:           req.AvailableTo,
	})
	if err != nil {
		return nil, fmt.Errorf("new resource: %w", err)
	}

	if err := c.resourceRepo.Create(ctx, res); err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	return &port.CreateResponse{ResourceID: res.ID()}, nil
}
