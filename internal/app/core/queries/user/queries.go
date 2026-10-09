package userqueries

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/user"
	userout "time4book/internal/app/core/ports/out/user"
)

func New(urepo userout.Repository, log *slog.Logger) port.Queries {
	return port.Queries{Get: newGet(urepo, log)}
}
