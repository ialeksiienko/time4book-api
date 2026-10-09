package reservationcommands

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

func New(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	brepo reservationout.Repository,
	v *validator.Facade,
	log *slog.Logger,
) port.Commands {
	return port.Commands{
		Cancel: newCancel(urepo, brepo, log.With(slog.String("reservation", "cancel"))),
		Create: newCreate(urepo, resrepo, brepo, v, log.With(slog.String("reservation", "create"))),
	}
}
