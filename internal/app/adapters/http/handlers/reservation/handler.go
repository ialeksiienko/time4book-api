package reservation

import (
	reservationcommands "time4book/internal/app/core/usecases/reservation"
)

type Handler struct {
	commands *reservationcommands.Facade
}

func NewHandler(commands *reservationcommands.Facade) *Handler {
	return &Handler{
		commands: commands,
	}
}
