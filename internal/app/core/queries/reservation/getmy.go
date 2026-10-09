package reservationqueries

import (
	"context"
	"fmt"
	"log/slog"
	"time4book/internal/app/core/domain/model/reservation"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
)

type GetMy struct {
	reservationRepo reservationout.Repository
	log             *slog.Logger
}

func newGetMy(
	brepo reservationout.Repository,
	l *slog.Logger,
) *GetMy {
	return &GetMy{
		reservationRepo: brepo,
		log:             l,
	}
}

func (c *GetMy) Execute(ctx context.Context, req *port.GetMyRequest) (*port.GetMyResponse, error) {
	res, total, err := c.reservationRepo.GetByUserID(ctx, req.UserID, req.CompanyID, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("list reservations by user: %w", err)
	}

	reservations := make([]reservation.Reservation, len(res))
	for i, r := range res {
		reservations[i] = *r
	}

	return &port.GetMyResponse{
		Reservations: reservations,
		Total:        total,
	}, nil
}
