package authcommands

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"
)

type Facade struct {
	Register *Register
	Login    *Login
	Me       *Me
	Logout   *Logout
	Refresh  *Refresh
}

func NewFacade(
	urepo user.UserRepo,
	arepo auth.AuthRepo,
	crepo company.CompanyRepo,
	tx ports.TxManager,
	token ports.TokenManager,
	v *validator.Facade,
	log *slog.Logger,
) *Facade {
	return &Facade{
		Register: newRegister(urepo, arepo, crepo, tx, token, v, log),
		Login:    newLogin(urepo, arepo, token, v, log),
		Me:       newMe(urepo, crepo, log),
		Logout:   newLogout(arepo, log),
		Refresh:  newRefresh(urepo, arepo, token, v, log),
	}
}
