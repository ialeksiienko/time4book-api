package httpadapter

import (
	authhandler "time4book/internal/app/adapters/in/http/handlers/auth"
	companyhandler "time4book/internal/app/adapters/in/http/handlers/company"
	companyresourcetypehandler "time4book/internal/app/adapters/in/http/handlers/companyresourcetype"
	reservationhandler "time4book/internal/app/adapters/in/http/handlers/reservation"
	resourcehandler "time4book/internal/app/adapters/in/http/handlers/resource"
	userhandler "time4book/internal/app/adapters/in/http/handlers/user"
	inports "time4book/internal/app/core/ports/in"
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
	p inports.Ports,
) *Handler {
	return &Handler{
		AuthHandler:                authhandler.NewHandler(p.Commands.Auth, p.Queries.Auth),
		UserHandler:                userhandler.NewHandler(p.Commands.User, p.Queries.User),
		CompanyHandler:             companyhandler.NewHandler(p.Commands.Company, p.Queries.Company),
		CompanyResourceTypeHandler: companyresourcetypehandler.NewHandler(p.Commands.CompanyResourceType, p.Queries.CompanyResourceType),
		ResourceHandler:            resourcehandler.NewHandler(p.Commands.Resource, p.Queries.Resource),
		ReservationHandler:         reservationhandler.NewHandler(p.Commands.Reservation, p.Queries.Reservation),
	}
}
