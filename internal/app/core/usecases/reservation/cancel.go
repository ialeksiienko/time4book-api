package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"
	"time4book/internal/app/core/domain/model/user"

	"github.com/google/uuid"
)

type CancelRequest struct {
	InitiatorID   uuid.UUID
	CompanyID     uuid.UUID
	ReservationID uuid.UUID
}

type CancelResponse struct{}

type Cancel struct {
	userRepo        user.UserRepo
	reservationRepo reservation.ReservationRepo
	log             *slog.Logger
}

func newCancel(
	urepo user.UserRepo,
	brepo reservation.ReservationRepo,
	l *slog.Logger,
) *Cancel {
	return &Cancel{
		userRepo:        urepo,
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *Cancel) Execute(ctx context.Context, req *CancelRequest) (*CancelResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	b, err := c.reservationRepo.ByID(ctx, req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("get reservation: %w", err)
	}

	if req.CompanyID != b.CompanyID() {
		return nil, fmt.Errorf("reservation not in company")
	}

	if b.UserID() == initiator.ID() {
		if err := b.Cancel(); err != nil {
			return nil, fmt.Errorf("cancel reservation: %w", err)
		}
	} else if initiator.Role().IsDeveloper() || initiator.Role().IsOwner() || initiator.Role().IsAdmin() {
		if err := b.CancelByAdmin(); err != nil {
			return nil, fmt.Errorf("cancel reservation by admin: %w", err)
		}
	} else {
		return nil, user.ErrUnauthorized
	}

	if err := c.reservationRepo.Update(ctx, b); err != nil {
		return nil, fmt.Errorf("update reservation: %w", err)
	}

	return &CancelResponse{}, nil
}
