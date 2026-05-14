package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"

	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID           uuid.UUID
	CompanyID             uuid.UUID
	Name                  string `validate:"required"`
	Type                  string `validate:"required"`
	CompanyResourceTypeID *uuid.UUID
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
}

type CreateResponse struct {
	ResourceID uuid.UUID
}

type Create struct {
	userRepo        user.UserRepo
	resourceRepo    resource.ResourceRepo
	companyTypeRepo companyresourcetype.Repo
	validator       *validator.Facade
	log             *slog.Logger
}

func newCreate(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	companyTypeRepo companyresourcetype.Repo,
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

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	rt := resource.ResourceType(req.Type)
	if rt != resource.TypeCustom {
		return nil, fmt.Errorf("only local company resource types are supported")
	}
	if req.CompanyResourceTypeID == nil {
		return nil, fmt.Errorf("company_resource_type_id is required when type is custom")
	}
	if _, err := c.companyTypeRepo.ByIDAndCompany(ctx, *req.CompanyResourceTypeID, req.CompanyID); err != nil {
		return nil, fmt.Errorf("company resource type: %w", err)
	}

	res, err := resource.NewResource(
		req.CompanyID,
		req.Name,
		rt,
		req.CompanyResourceTypeID,
		req.Description,
		req.Location,
		req.MaxReservationMinutes,
		req.AvailableFrom,
		req.AvailableTo,
	)
	if err != nil {
		return nil, fmt.Errorf("new resource: %w", err)
	}

	if err := c.resourceRepo.Create(ctx, res); err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	return &CreateResponse{ResourceID: res.ID()}, nil
}
