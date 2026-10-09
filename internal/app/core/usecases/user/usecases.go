package usercommands

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/user"
	outports "time4book/internal/app/core/ports/out"
	authout "time4book/internal/app/core/ports/out/auth"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

func New(
	urepo userout.Repository,
	arepo authout.Repository,
	tx outports.TxManager,
	v *validator.Facade,
	log *slog.Logger,
) port.Commands {
	return port.Commands{
		Create:     newCreate(urepo, arepo, tx, v, log),
		Deactivate: newDeactivate(urepo, log),
		Update:     newUpdate(urepo, v, log),
	}
}
