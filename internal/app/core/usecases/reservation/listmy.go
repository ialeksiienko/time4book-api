package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

type ListMyRequest struct {
	UserID uuid.UUID
	Page   int
	Limit  int
}

type ListMyResponse struct {
	Reservations []booking.Booking
	Total        int64
}

type ListMy struct {
	bookingRepo booking.BookingRepo
	log         *slog.Logger
}

func newListMy(
	brepo booking.BookingRepo,
	l *slog.Logger,
) *ListMy {
	return &ListMy{
		bookingRepo: brepo,
		log:         l,
	}
}

func (c *ListMy) Execute(ctx context.Context, req *ListMyRequest) (*ListMyResponse, error) {
	res, total, err := c.bookingRepo.ListByUserID(ctx, req.UserID, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by user: %w", err)
	}

	reservations := make([]booking.Booking, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &ListMyResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
