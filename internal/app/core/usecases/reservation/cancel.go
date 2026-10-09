package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
	userout "time4book/internal/app/core/ports/out/user"
)

type Cancel struct {
	userRepo        userout.Repository
	reservationRepo reservationout.Repository
	log             *slog.Logger
}

func newCancel(
	urepo userout.Repository,
	brepo reservationout.Repository,
	l *slog.Logger,
) *Cancel {
	return &Cancel{
		userRepo:        urepo,
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *Cancel) Execute(ctx context.Context, req *port.CancelRequest) (*port.CancelResponse, error) {
	initiator, err := c.userRepo.Retrieve(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	b, err := c.reservationRepo.Retrieve(ctx, req.ReservationID)
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

	return &port.CancelResponse{}, nil
}
