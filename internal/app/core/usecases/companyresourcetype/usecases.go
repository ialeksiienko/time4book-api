package companyresourcetypecommands

import (
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	userout "time4book/internal/app/core/ports/out/user"
)

func New(repo companyresourcetypeout.Repository, userRepo userout.Repository) port.Commands {
	return port.Commands{
		Create: &Create{repo: repo, userRepo: userRepo},
		Delete: &Delete{repo: repo, userRepo: userRepo},
		Update: &Update{repo: repo, userRepo: userRepo},
	}
}
