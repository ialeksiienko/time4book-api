package resourcecommands

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"
)

type Facade struct {
	Create   *Create
	List     *List
	GetByID  *GetByID
	Update   *Update
	Delete   *Delete
	Service  *Service
	Restore  *Restore
}

func NewFacade(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	v *validator.Facade,
	log *slog.Logger,
) *Facade {
	return &Facade{
		Create:   newCreate(urepo, resrepo, v, log.With(slog.String("resource", "create"))),
		List:     newList(resrepo, log.With(slog.String("resource", "list"))),
		GetByID:  newGetByID(resrepo, log.With(slog.String("resource", "getByID"))),
		Update:   newUpdate(urepo, resrepo, v, log.With(slog.String("resource", "update"))),
		Delete:   newDelete(urepo, resrepo, log.With(slog.String("resource", "delete"))),
		Service:  newService(urepo, resrepo, v, log.With(slog.String("resource", "service"))),
		Restore:  newRestore(urepo, resrepo, log.With(slog.String("resource", "restore"))),
	}
}
