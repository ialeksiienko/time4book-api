package user

import (
	usercommands "time4book/internal/app/core/usecases/user"
)

type Handler struct {
	commands *usercommands.Facade
}

func NewHandler(commands *usercommands.Facade) *Handler {
	return &Handler{
		commands: commands,
	}
}
