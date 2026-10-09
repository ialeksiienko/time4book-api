package httpadapter

import (
	authhandler "time4book/internal/app/adapters/http/handlers/auth"
	companyhandler "time4book/internal/app/adapters/http/handlers/company"
	companyresourcetypehandler "time4book/internal/app/adapters/http/handlers/companyresourcetype"
	reservationhandler "time4book/internal/app/adapters/http/handlers/reservation"
	resourcehandler "time4book/internal/app/adapters/http/handlers/resource"
	userhandler "time4book/internal/app/adapters/http/handlers/user"
	"time4book/internal/app/core/usecases"
)

type Handler struct {
	AuthHandler                *authhandler.Handler
	UserHandler                *userhandler.Handler
	CompanyHandler             *companyhandler.Handler
	CompanyResourceTypeHandler *companyresourcetypehandler.Handler
	ResourceHandler            *resourcehandler.Handler
	ReservationHandler         *reservationhandler.Handler
}

func NewHandler(
	c *usecases.Commands,
) *Handler {
	return &Handler{
		AuthHandler:                authhandler.NewHandler(c.Auth),
		UserHandler:                userhandler.NewHandler(c.User),
		CompanyHandler:             companyhandler.NewHandler(c.Company),
		CompanyResourceTypeHandler: companyresourcetypehandler.NewHandler(c.CompanyResourceType),
		ResourceHandler:            resourcehandler.NewHandler(c.Resource),
		ReservationHandler:         reservationhandler.NewHandler(c.Reservation),
	}
}
