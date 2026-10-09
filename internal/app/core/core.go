package core

import (
	"log/slog"
	inports "time4book/internal/app/core/ports/in"
	outports "time4book/internal/app/core/ports/out"
	authout "time4book/internal/app/core/ports/out/auth"
	companyout "time4book/internal/app/core/ports/out/company"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	reservationout "time4book/internal/app/core/ports/out/reservation"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	authqueries "time4book/internal/app/core/queries/auth"
	companyqueries "time4book/internal/app/core/queries/company"
	companyresourcetypequeries "time4book/internal/app/core/queries/companyresourcetype"
	reservationqueries "time4book/internal/app/core/queries/reservation"
	resourcequeries "time4book/internal/app/core/queries/resource"
	userqueries "time4book/internal/app/core/queries/user"
	authcommands "time4book/internal/app/core/usecases/auth"
	companycommands "time4book/internal/app/core/usecases/company"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"
	reservationcommands "time4book/internal/app/core/usecases/reservation"
	resourcecommands "time4book/internal/app/core/usecases/resource"
	usercommands "time4book/internal/app/core/usecases/user"
	"time4book/pkg/validator"
)

func New(
	userRepo userout.Repository,
	authRepo authout.Repository,
	companyRepo companyout.Repository,
	companyTypeRepo companyresourcetypeout.Repository,
	resourceRepo resourceout.Repository,
	resourceReader resourceout.Reader,
	reservationRepo reservationout.Repository,
	txManager outports.TxManager,
	tokenManager outports.TokenManager,
	validator *validator.Facade,
	log *slog.Logger,
) inports.Ports {
	return inports.Ports{
		Commands: inports.Commands{
			Auth:                authcommands.New(userRepo, authRepo, companyRepo, txManager, tokenManager, validator, log),
			Company:             companycommands.New(userRepo, companyRepo, txManager, validator, log),
			CompanyResourceType: companyresourcetypecommands.New(companyTypeRepo, userRepo),
			Reservation:         reservationcommands.New(userRepo, resourceRepo, reservationRepo, validator, log),
			Resource:            resourcecommands.New(userRepo, resourceRepo, companyTypeRepo, validator, log),
			User:                usercommands.New(userRepo, authRepo, txManager, validator, log),
		},
		Queries: inports.Queries{
			Auth:                authqueries.New(userRepo, companyRepo, log),
			Company:             companyqueries.New(companyRepo, log),
			CompanyResourceType: companyresourcetypequeries.New(companyTypeRepo, userRepo),
			Reservation:         reservationqueries.New(reservationRepo, log),
			Resource:            resourcequeries.New(resourceReader, log),
			User:                userqueries.New(userRepo, log),
		},
	}
}
