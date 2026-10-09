package company

import (
	port "time4book/internal/app/core/ports/in/company"
)

type Handler struct {
	commands port.Commands
	queries  port.Queries
}

func NewHandler(commands port.Commands, queries port.Queries) *Handler {
	return &Handler{
		commands: commands, queries: queries,
	}
}
