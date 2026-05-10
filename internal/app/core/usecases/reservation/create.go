package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/reservation"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID uuid.UUID `validate:"required"`
	CompanyID   uuid.UUID `validate:"required"`
	ResourceID  uuid.UUID `validate:"required"`
	StartDate   time.Time `validate:"required"`
	EndDate     time.Time `validate:"required"`
	Description *string
}

type CreateResponse struct {
	ReservationID uuid.UUID
}

type Create struct {
	userRepo        user.UserRepo
	resourceRepo    resource.ResourceRepo
	reservationRepo reservation.ReservationRepo
	validator       *validator.Facade
	log             *slog.Logger
}

func newCreate(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	brepo reservation.ReservationRepo,
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

func (c *Create) Execute(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
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

	if initiator.CompanyID() == uuid.Nil || initiator.CompanyID() != res.CompanyID() {
		return nil, user.ErrUnauthorized
	}

	if !res.IsBookable() {
		return nil, fmt.Errorf("resource is not bookable")
	}

	if res.MaxReservationMinutes() != nil {
		dur := req.EndDate.Sub(req.StartDate)
		if int(dur.Minutes()) > *res.MaxReservationMinutes() {
			return nil, fmt.Errorf("reservation exceeds maximum allowed time")
		}
	}

	activeReservations, err := c.reservationRepo.ActiveByResourceIDInRange(ctx, req.ResourceID, req.CompanyID, req.StartDate, req.EndDate, nil)
	if err != nil {
		return nil, fmt.Errorf("check active reservations: %w", err)
	}
	if len(activeReservations) > 0 {
		return nil, fmt.Errorf("time slot is already booked")
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

	return &CreateResponse{ReservationID: b.ID()}, nil
}
