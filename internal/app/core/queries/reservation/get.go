package reservationqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
)

type Get struct {
	reservationRepo reservationout.Repository
	log             *slog.Logger
}

func newGet(
	brepo reservationout.Repository,
	l *slog.Logger,
) *Get {
	return &Get{
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *Get) Execute(ctx context.Context, req *port.GetRequest) (*port.GetResponse, error) {
	filter := reservationout.Filter{
		CompanyID:  req.CompanyID,
		UserID:     req.UserID,
		ResourceID: req.ResourceID,
		From:       req.From,
		To:         req.To,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	if req.Status != nil && *req.Status != "" {
		s := reservation.ReservationStatus(*req.Status)
		filter.Status = &s
	}

	res, total, err := c.reservationRepo.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}

	reservations := make([]reservation.Reservation, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &port.GetResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
