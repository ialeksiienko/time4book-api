package httpadapter

import (
	authhandler "time4book/internal/app/adapters/in/http/handlers/auth"
	companyhandler "time4book/internal/app/adapters/in/http/handlers/company"
	reservationhandler "time4book/internal/app/adapters/in/http/handlers/reservation"
	resourcehandler "time4book/internal/app/adapters/in/http/handlers/resource"
	userhandler "time4book/internal/app/adapters/in/http/handlers/user"
	"time4book/internal/app/core/usecases"
)

type Handler struct {
	AuthHandler        *authhandler.Handler
	UserHandler        *userhandler.Handler
	CompanyHandler     *companyhandler.Handler
	ResourceHandler    *resourcehandler.Handler
	ReservationHandler *reservationhandler.Handler
}

func NewHandler(
	c *usecases.Commands,
) *Handler {
	return &Handler{
		AuthHandler:        authhandler.NewHandler(c.Auth),
		UserHandler:        userhandler.NewHandler(c.User),
		CompanyHandler:     companyhandler.NewHandler(c.Company),
		ResourceHandler:    resourcehandler.NewHandler(c.Resource),
		ReservationHandler: reservationhandler.NewHandler(c.Reservation),
	}
}

