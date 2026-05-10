package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
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
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
}

type UpdateResponse struct{}

type Update struct {
	userRepo     user.UserRepo
	resourceRepo resource.ResourceRepo
	validator    *validator.Facade
	log          *slog.Logger
}

func newUpdate(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Update {
	return &Update{
		userRepo:     urepo,
		resourceRepo: resrepo,
		validator:    v,
		log:          l,
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

	res, err := c.resourceRepo.ByID(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != res.CompanyID() {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() && !initiator.Role().IsAdmin() {
			return nil, user.ErrUnauthorized
		}
	}

	props := &resource.Props{
		ID:                    res.ID(),
		CompanyID:             res.CompanyID(),
		Name:                  req.Name,
		ResourceType:          req.Type,
		Description:           req.Description,
		Location:              req.Location,
		MaxReservationMinutes: req.MaxReservationMinutes,
		AvailableFrom:         req.AvailableFrom,
		AvailableTo:           req.AvailableTo,
		Status:                res.Status().String(),
		UnavailableFrom:       res.UnavailableFrom(),
		UnavailableTo:         res.UnavailableTo(),
		UnavailableReason:     res.UnavailableReason(),
		CreatedAt:             res.CreatedAt(),
	}

	updatedRes := resource.Reconstitute(props)
	updatedRes.Restore()
	if res.Status() == resource.StatusInService {
		updatedRes.MarkInService(*res.UnavailableReason(), *res.UnavailableFrom(), res.UnavailableTo())
	} else if res.Status() == resource.StatusInactive {
		updatedRes.Deactivate()
	}

	if err := c.resourceRepo.Update(ctx, updatedRes); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &UpdateResponse{}, nil
}
