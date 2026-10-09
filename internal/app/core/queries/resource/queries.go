package resourcequeries

import (
	"log/slog"
	port "time4book/internal/app/core/ports/in/resource"
	resourceout "time4book/internal/app/core/ports/out/resource"
)

func New(reader resourceout.Reader, log *slog.Logger) port.Queries {
	return port.Queries{
		Retrieve: newRetrieve(reader, log.With(slog.String("resource", "getByID"))),
		Get:      newGet(reader, log.With(slog.String("resource", "list"))),
	}
}
