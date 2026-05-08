package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

type ListByResourceRequest struct {
	ResourceID uuid.UUID
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type ListByResourceResponse struct {
	Reservations []booking.Booking
	Total        int64
}

type ListByResource struct {
	bookingRepo booking.BookingRepo
	log         *slog.Logger
}

func newListByResource(
	brepo booking.BookingRepo,
	l *slog.Logger,
) *ListByResource {
	return &ListByResource{
		bookingRepo: brepo,
		log:         l,
	}
}

func (c *ListByResource) Execute(ctx context.Context, req *ListByResourceRequest) (*ListByResourceResponse, error) {
	res, total, err := c.bookingRepo.ListByResourceID(ctx, req.ResourceID, req.From, req.To, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by resource: %w", err)
	}

	reservations := make([]booking.Booking, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &ListByResourceResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
