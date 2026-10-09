package reservation

import (
	port "time4book/internal/app/core/ports/in/reservation"
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
