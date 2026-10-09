package reservationqueries

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
)

func New(brepo reservationout.Repository, log *slog.Logger) port.Queries {
	return port.Queries{
		Get:           newGet(brepo, log.With(slog.String("reservation", "list"))),
		GetByResource: newGetByResource(brepo, log.With(slog.String("reservation", "listByResource"))),
		GetMy:         newGetMy(brepo, log.With(slog.String("reservation", "listMy"))),
	}
}
