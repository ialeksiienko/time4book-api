package resourcecommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type ServiceRequest struct {
	InitiatorID uuid.UUID
	ResourceID  uuid.UUID
	Reason      string    `validate:"required"`
	From        time.Time `validate:"required"`
	To          *time.Time
}

type ServiceResponse struct{}

type Service struct {
	userRepo     user.UserRepo
	resourceRepo resource.ResourceRepo
	validator    *validator.Facade
	log          *slog.Logger
}

func newService(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
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

func (c *Service) Execute(ctx context.Context, req *ServiceRequest) (*ServiceResponse, error) {
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

	res.MarkInService(req.Reason, req.From, req.To)

	if err := c.resourceRepo.Update(ctx, res); err != nil {
		return nil, fmt.Errorf("update resource: %w", err)
	}

	return &ServiceResponse{}, nil
}
