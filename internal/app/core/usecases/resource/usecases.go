package resourcecommands

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/resource"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"
)

func New(
	urepo userout.Repository,
	resrepo resourceout.Repository,
	companyTypeRepo companyresourcetypeout.Repository,
	v *validator.Facade,
	log *slog.Logger,
) port.Commands {
	return port.Commands{
		Create:  newCreate(urepo, resrepo, companyTypeRepo, v, log.With(slog.String("resource", "create"))),
		Delete:  newDelete(urepo, resrepo, log.With(slog.String("resource", "delete"))),
		Restore: newRestore(urepo, resrepo, log.With(slog.String("resource", "restore"))),
		Service: newService(urepo, resrepo, v, log.With(slog.String("resource", "service"))),
		Update:  newUpdate(urepo, resrepo, companyTypeRepo, v, log.With(slog.String("resource", "update"))),
	}
}
