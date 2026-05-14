package usecases

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/auth"
	"time4book/internal/app/core/domain/model/booking"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	authcommands "time4book/internal/app/core/usecases/auth"
	companycommands "time4book/internal/app/core/usecases/company"
	companyresourcetypecommands "time4book/internal/app/core/usecases/companyresourcetype"
	reservationcommands "time4book/internal/app/core/usecases/reservation"
	resourcecommands "time4book/internal/app/core/usecases/resource"
	usercommands "time4book/internal/app/core/usecases/user"
	"time4book/pkg/validator"
)

type Commands struct {
	Auth                *authcommands.Facade
	User                *usercommands.Facade
	Company             *companycommands.Facade
	Resource            *resourcecommands.Facade
	Reservation         *reservationcommands.Facade
	CompanyResourceType *companyresourcetypecommands.Facade
}

func New(
	userRepo user.UserRepo,
	authRepo auth.AuthRepo,
	companyRepo company.CompanyRepo,
	resourceRepo resource.ResourceRepo,
	companyResourceTypes companyresourcetype.Repo,
	bookingRepo booking.BookingRepo,
	txManager ports.TxManager,
	tokenManager ports.TokenManager,
	validator *validator.Facade,
	log *slog.Logger,
) *Commands {
	return &Commands{
		Auth:                authcommands.NewFacade(userRepo, authRepo, companyRepo, txManager, tokenManager, validator, log.With(slog.String("module", "auth"))),
		User:                usercommands.NewFacade(userRepo, authRepo, txManager, validator, log.With(slog.String("module", "user"))),
		Company:             companycommands.NewFacade(userRepo, companyRepo, bookingRepo, txManager, validator, log.With(slog.String("module", "company"))),
		Resource:            resourcecommands.NewFacade(userRepo, resourceRepo, companyResourceTypes, validator, log.With(slog.String("module", "resource"))),
		Reservation:         reservationcommands.NewFacade(userRepo, resourceRepo, bookingRepo, validator, log.With(slog.String("module", "reservation"))),
		CompanyResourceType: companyresourcetypecommands.NewFacade(companyResourceTypes, userRepo, validator, log.With(slog.String("module", "company_resource_type"))),
	}
}
