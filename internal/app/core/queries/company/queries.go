package companyqueries

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/company"
	companyout "time4book/internal/app/core/ports/out/company"
)

func New(crepo companyout.Repository, log *slog.Logger) port.Queries {
	return port.Queries{
		Retrieve: newRetrieve(crepo, log.With(slog.String("company", "getByID"))),
		Get:      newGet(crepo, log.With(slog.String("company", "list"))),
	}
}
