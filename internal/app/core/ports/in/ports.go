package inports

import (
	authport "time4book/internal/app/core/ports/in/auth"
	companyport "time4book/internal/app/core/ports/in/company"
	companyresourcetypeport "time4book/internal/app/core/ports/in/companyresourcetype"
	reservationport "time4book/internal/app/core/ports/in/reservation"
	resourceport "time4book/internal/app/core/ports/in/resource"
	userport "time4book/internal/app/core/ports/in/user"
)

type Commands struct {
	Auth                authport.Commands
	Company             companyport.Commands
	CompanyResourceType companyresourcetypeport.Commands
	Reservation         reservationport.Commands
	Resource            resourceport.Commands
	User                userport.Commands
}

type Queries struct {
	Auth                authport.Queries
	Company             companyport.Queries
	CompanyResourceType companyresourcetypeport.Queries
	Reservation         reservationport.Queries
	Resource            resourceport.Queries
	User                userport.Queries
}

type Ports struct {
	Commands Commands
	Queries  Queries
}
