package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

type ListMyRequest struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Page      int
	Limit     int
}

type ListMyResponse struct {
	Reservations []reservation.Reservation
	Total        int64
}

type ListMy struct {
	reservationRepo reservation.ReservationRepo
	log             *slog.Logger
}

func newListMy(
	brepo reservation.ReservationRepo,
	l *slog.Logger,
) *ListMy {
	return &ListMy{
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *ListMy) Execute(ctx context.Context, req *ListMyRequest) (*ListMyResponse, error) {
	res, total, err := c.reservationRepo.ListByUserID(ctx, req.UserID, req.CompanyID, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by user: %w", err)
	}

	reservations := make([]reservation.Reservation, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &ListMyResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
