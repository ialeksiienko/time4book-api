package authqueries

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/auth"
	companyout "time4book/internal/app/core/ports/out/company"
	userout "time4book/internal/app/core/ports/out/user"
)

func New(urepo userout.Repository, crepo companyout.Repository, log *slog.Logger) port.Queries {
	return port.Queries{RetrieveProfile: newRetrieveProfile(urepo, crepo, log)}
}
