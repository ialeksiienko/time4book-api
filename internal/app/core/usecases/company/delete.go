package companycommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/booking"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"

	"github.com/google/uuid"
)

type DeleteRequest struct {
	InitiatorID uuid.UUID
	CompanyID   uuid.UUID
}

type DeleteResponse struct{}

type Delete struct {
	userRepo    user.UserRepo
	companyRepo company.CompanyRepo
	bookingRepo booking.BookingRepo
	tx          ports.TxManager
	log         *slog.Logger
}

func newDelete(
	urepo user.UserRepo,
	crepo company.CompanyRepo,
	brepo booking.BookingRepo,
	tx ports.TxManager,
	l *slog.Logger,
) *Delete {
	return &Delete{
		userRepo:    urepo,
		companyRepo: crepo,
		bookingRepo: brepo,
		tx:          tx,
		log:         l,
	}
}

func (c *Delete) Execute(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	initiator, err := c.userRepo.ByID(ctx, req.InitiatorID)
	if err != nil {
		return nil, fmt.Errorf("get initiator: %w", err)
	}

	comp, err := c.companyRepo.ByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	if !initiator.Role().IsDeveloper() {
		if initiator.CompanyID() == nil || *initiator.CompanyID() != req.CompanyID {
			return nil, user.ErrUnauthorized
		}
		if !initiator.Role().IsOwner() || initiator.ID() != comp.OwnerID() {
			return nil, user.ErrUnauthorized
		}
	}

	if err := c.tx.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := c.bookingRepo.DeleteByCompanyID(txCtx, req.CompanyID); err != nil {
			return fmt.Errorf("delete company reservations: %w", err)
		}
		if err := c.companyRepo.Delete(txCtx, req.CompanyID); err != nil {
			return fmt.Errorf("delete company: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &DeleteResponse{}, nil
}
