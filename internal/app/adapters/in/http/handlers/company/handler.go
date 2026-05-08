package company

import (
	companycommands "time4book/internal/app/core/usecases/company"
)

type Handler struct {
	commands *companycommands.Facade
}

func NewHandler(commands *companycommands.Facade) *Handler {
	return &Handler{
		commands: commands,
	}
}
