package companyresourcetypecommands

import (
	"log/slog"

	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"
)

type Facade struct {
	Create *Create
	List   *List
	Update *Update
	Delete *Delete
}

func NewFacade(repo companyresourcetype.Repo, ur user.UserRepo, v *validator.Facade, log *slog.Logger) *Facade {
	return &Facade{
		Create: newCreate(repo, ur, v, log.With(slog.String("module", "company_resource_type")).With(slog.String("usecase", "create"))),
		List:   newList(repo, ur, log.With(slog.String("module", "company_resource_type")).With(slog.String("usecase", "list"))),
		Update: newUpdate(repo, ur, v, log.With(slog.String("module", "company_resource_type")).With(slog.String("usecase", "update"))),
		Delete: newDelete(repo, ur, log.With(slog.String("module", "company_resource_type")).With(slog.String("usecase", "delete"))),
	}
}
