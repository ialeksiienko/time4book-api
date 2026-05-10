package reservationcommands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"time4book/internal/app/core/domain/model/reservation"

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
	Reservations []reservation.Reservation
	Total        int64
}

type ListByResource struct {
	reservationRepo reservation.ReservationRepo
	log             *slog.Logger
}

func newListByResource(
	brepo reservation.ReservationRepo,
	l *slog.Logger,
) *ListByResource {
	return &ListByResource{
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *ListByResource) Execute(ctx context.Context, req *ListByResourceRequest) (*ListByResourceResponse, error) {
	res, total, err := c.reservationRepo.ListByResourceID(ctx, req.ResourceID, req.From, req.To, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by resource: %w", err)
	}

	reservations := make([]reservation.Reservation, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &ListByResourceResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
