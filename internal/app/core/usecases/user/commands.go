package usercommands

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"
)

type Facade struct {
	List       *List
	Create     *Create
	Update     *Update
	Deactivate *Deactivate
}

func NewFacade(
	urepo user.UserRepo,
	arepo auth.AuthRepo,
	tx ports.TxManager,
	v *validator.Facade,
	log *slog.Logger,
) *Facade {
	return &Facade{
		List:       newList(urepo, log),
		Create:     newCreate(urepo, arepo, tx, v, log),
		Update:     newUpdate(urepo, v, log),
		Deactivate: newDeactivate(urepo, log),
	}
}
