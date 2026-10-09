package authcommands

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/auth"
	outports "time4book/internal/app/core/ports/out"
	authout "time4book/internal/app/core/ports/out/auth"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

func New(
	urepo userout.Repository,
	arepo authout.Repository,
	crepo companyout.Repository,
	tx outports.TxManager,
	token outports.TokenManager,
	v *validator.Facade,
	log *slog.Logger,
) port.Commands {
	return port.Commands{
		Login:    newLogin(urepo, arepo, token, v, log),
		Logout:   newLogout(arepo, log),
		Refresh:  newRefresh(urepo, arepo, token, v, log),
		Register: newRegister(urepo, arepo, crepo, tx, token, v, log),
	}
}
