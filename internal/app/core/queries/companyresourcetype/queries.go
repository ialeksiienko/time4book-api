package companyresourcetypequeries

import (
	port "time4book/internal/app/core/ports/in/companyresourcetype"
	companyresourcetypeout "time4book/internal/app/core/ports/out/companyresourcetype"
	userout "time4book/internal/app/core/ports/out/user"
)

func New(repo companyresourcetypeout.Repository, userRepo userout.Repository) port.Queries {
	return port.Queries{Get: &Get{repo: repo, userRepo: userRepo}}
}
