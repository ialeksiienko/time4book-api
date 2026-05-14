package companyresourcetype

import (
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"
)

type Handler struct {
	commands *companyresourcetypecommands.Facade
}

func NewHandler(f *companyresourcetypecommands.Facade) *Handler {
	return &Handler{commands: f}
}
