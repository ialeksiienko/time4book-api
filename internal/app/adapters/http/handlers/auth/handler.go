package auth

import (
	authcommands "time4book/internal/app/core/usecases/auth"
)

type Handler struct {
	commands *authcommands.Facade
}

func NewHandler(commands *authcommands.Facade) *Handler {
	return &Handler{
		commands: commands,
	}
}

