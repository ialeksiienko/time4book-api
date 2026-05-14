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

type UpdateRequest struct {
	InitiatorID           uuid.UUID
	ResourceID            uuid.UUID
	Name                  string `validate:"required"`
	Type                  string `validate:"required"`
	CompanyResourceTypeID *uuid.UUID
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
}

type UpdateResponse struct{}

type Update struct {
	userRepo        user.UserRepo
	resourceRepo    resource.ResourceRepo
	companyTypeRepo companyresourcetype.Repo
	validator       *validator.Facade
	log             *slog.Logger
}

func newUpdate(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	companyTypeRepo companyresourcetype.Repo,
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

func (c *Update) Execute(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validate error: %w", err)
	}

	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	old, err := c.resourceRepo.ByID(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != old.CompanyID() {
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
	if _, err := c.companyTypeRepo.ByIDAndCompany(ctx, *req.CompanyResourceTypeID, old.CompanyID()); err != nil {
		return nil, fmt.Errorf("company resource type: %w", err)
	}

	props := &resource.Props{
		ID:                    old.ID(),
		CompanyID:             old.CompanyID(),
		Name:                  req.Name,
		ResourceType:          req.Type,
		CompanyResourceTypeID: req.CompanyResourceTypeID,
		Description:           req.Description,
		Location:              req.Location,
		MaxReservationMinutes: req.MaxReservationMinutes,
		AvailableFrom:         req.AvailableFrom,
		AvailableTo:           req.AvailableTo,
		Status:                old.Status().String(),
		UnavailableFrom:       old.UnavailableFrom(),
		UnavailableTo:         old.UnavailableTo(),
		UnavailableReason:     old.UnavailableReason(),
		CreatedAt:             old.CreatedAt(),
		UpdatedAt:             old.UpdatedAt(),
	}

	updatedRes := resource.Reconstitute(props)

	updatedRes.Restore()

	if old.Status() == resource.StatusInService {
		updatedRes.MarkInService(*old.UnavailableReason(), *old.UnavailableFrom(), old.UnavailableTo())
	} else if old.Status() == resource.StatusInactive {
		updatedRes.Deactivate()
	}

	if err := c.resourceRepo.Update(ctx, updatedRes); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &UpdateResponse{}, nil
}
