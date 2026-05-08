package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type UpdateRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
	Name        string `validate:"required"`
	NIP         *string
	Address     *string
	Industry    *string
}

type UpdateResponse struct{}

type Update struct {
	userRepo    user.UserRepo
	companyRepo company.CompanyRepo
	validator   *validator.Facade
	log         *slog.Logger
}

func newUpdate(
	urepo user.UserRepo,
	crepo company.CompanyRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Update {
	return &Update{
		userRepo:    urepo,
		companyRepo: crepo,
		validator:   v,
		log:         l,
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

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() {
			return nil, user.ErrUnauthorized
		}
	}

	comp, err := c.companyRepo.ByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	props := &company.Props{
		ID:        comp.ID(),
		OwnerID:   comp.OwnerID(),
		Name:      req.Name,
		NIP:       req.NIP,
		Address:   req.Address,
		Industry:  req.Industry,
		Status:    comp.Status().String(),
		CreatedAt: comp.CreatedAt(),
	}
	updatedComp := company.Reconstitute(props)
	updatedComp.Unblock()
	if comp.IsBlocked() {
		updatedComp.Block()
	}

	if err := c.companyRepo.Update(ctx, updatedComp); err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}

	return &UpdateResponse{}, nil
}
