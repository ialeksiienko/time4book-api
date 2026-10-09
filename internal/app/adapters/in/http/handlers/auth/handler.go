package auth

import (
	port "time4book/internal/app/core/ports/in/auth"
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
