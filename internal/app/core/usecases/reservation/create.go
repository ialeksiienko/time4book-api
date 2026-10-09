package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type Create struct {
	userRepo        userout.Repository
	resourceRepo    resourceout.Repository
	reservationRepo reservationout.Repository
	validator       *validator.Facade
	log             *slog.Logger
}

func newCreate(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	brepo reservationout.Repository,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{
		userRepo:        urepo,
		resourceRepo:    resrepo,
		reservationRepo: brepo,
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

	res, err := c.resourceRepo.Retrieve(ctx, req.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get resource: %w", err)
	}

	if !initiator.IsActive() || initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != req.CompanyID || req.CompanyID != res.CompanyID() {
		return nil, user.ErrUnauthorized
	}

	if err := res.ValidateBookingInterval(req.StartDate, req.EndDate); err != nil {
		return nil, fmt.Errorf("validate booking interval: %w", err)
	}

	b, err := reservation.NewReservation(
		initiator.ID(),
		req.CompanyID,
		req.ResourceID,
		req.StartDate,
		req.EndDate,
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation: %w", err)
	}

	if err := c.reservationRepo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create reservation: %w", err)
	}

	return &port.CreateResponse{ReservationID: b.ID()}, nil
}
