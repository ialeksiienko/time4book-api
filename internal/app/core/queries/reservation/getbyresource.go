package reservationqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
)

type GetByResource struct {
	reservationRepo reservationout.Repository
	log             *slog.Logger
}

func newGetByResource(
	brepo reservationout.Repository,
	l *slog.Logger,
) *GetByResource {
	return &GetByResource{
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *GetByResource) Execute(ctx context.Context, req *port.GetByResourceRequest) (*port.GetByResourceResponse, error) {
	res, total, err := c.reservationRepo.GetByResourceID(ctx, req.ResourceID, req.CompanyID, req.From, req.To, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by resource: %w", err)
	}

	reservations := make([]reservation.Reservation, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &port.GetByResourceResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
