package companycommands

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/company"
	outports "time4book/internal/app/core/ports/out"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

func New(
	urepo userout.Repository,
	crepo companyout.Repository,
	tx outports.TxManager,
	v *validator.Facade,
	log *slog.Logger,
) port.Commands {
	return port.Commands{
		Block:  newBlock(urepo, crepo, log.With(slog.String("company", "block"))),
		Create: newCreate(crepo, tx, v, log.With(slog.String("company", "create"))),
		Update: newUpdate(urepo, crepo, v, log.With(slog.String("company", "update"))),
	}
}
