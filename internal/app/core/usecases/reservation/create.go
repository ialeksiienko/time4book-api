package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/booking"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type CreateRequest struct {
	InitiatorID uuid.UUID
	ResourceID  uuid.UUID
	StartDate   time.Time `validate:"required"`
	EndDate     time.Time `validate:"required"`
	Description *string
}

type CreateResponse struct {
	ReservationID uuid.UUID
}

type Create struct {
	userRepo     user.UserRepo
	resourceRepo resource.ResourceRepo
	bookingRepo  booking.BookingRepo
	validator    *validator.Facade
	log          *slog.Logger
}

func newCreate(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	brepo booking.BookingRepo,
	v *validator.Facade,
	l *slog.Logger,
) *Create {
	return &Create{
		userRepo:     urepo,
		resourceRepo: resrepo,
		bookingRepo:  brepo,
		validator:    v,
		log:          l,
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

	if initiator.CompanyID() == nil || *initiator.CompanyID() != res.CompanyID() {
		return nil, user.ErrUnauthorized
	}

	if req.StartDate.Before(time.Now().UTC().Add(-time.Hour)) {
		return nil, ErrReservationStartTooFarInPast
	}

	if !res.IsBookableForInterval(req.StartDate, req.EndDate) {
		return nil, fmt.Errorf("zasób jest niedostępny w wybranym terminie")
	}

	if res.MaxReservationMinutes() != nil {
		dur := req.EndDate.Sub(req.StartDate)
		if int(dur.Minutes()) > *res.MaxReservationMinutes() {
			return nil, fmt.Errorf("reservation exceeds maximum allowed time")
		}
	}

	activeBookings, err := c.bookingRepo.ActiveByResourceIDInRange(ctx, req.ResourceID, req.StartDate, req.EndDate, nil)
	if err != nil {
		return nil, fmt.Errorf("check active bookings: %w", err)
	}
	if len(activeBookings) > 0 {
		return nil, ErrSlotAlreadyTaken
	}

	b, err := booking.NewBooking(
		initiator.ID(),
		req.ResourceID,
		req.StartDate,
		req.EndDate,
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("new booking: %w", err)
	}

	if err := c.bookingRepo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return &CreateResponse{ReservationID: b.ID()}, nil
}
