package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/booking"

	"github.com/google/uuid"
)

type ListRequest struct {
	CompanyID  *uuid.UUID
	UserID     *uuid.UUID
	ResourceID *uuid.UUID
	Status     *string
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type ListResponse struct {
	Reservations []booking.Booking
	Total        int64
}

type List struct {
	bookingRepo booking.BookingRepo
	log         *slog.Logger
}

func newList(
	brepo booking.BookingRepo,
	l *slog.Logger,
) *List {
	return &List{
		bookingRepo: brepo,
		log:         l,
	}
}

func (c *List) Execute(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	filter := booking.ListFilter{
		CompanyID:  req.CompanyID,
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		From:       req.From,
		To:         req.To,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	if req.Status != nil && *req.Status != "" {
		s := booking.BookingStatus(*req.Status)
		filter.Status = &s
	}

	res, total, err := c.bookingRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}

	reservations := make([]booking.Booking, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &ListResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
