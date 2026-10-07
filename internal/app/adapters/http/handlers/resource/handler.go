package resource

import (
	resourcecommands "time4book/internal/app/core/usecases/resource"
)

type Handler struct {
	commands *resourcecommands.Facade
}

func NewHandler(commands *resourcecommands.Facade) *Handler {
	return &Handler{
		commands: commands,
	}
}
