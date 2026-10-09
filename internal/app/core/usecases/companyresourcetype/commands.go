package companyresourcetypecommands

import (
	"time4book/internal/app/core/domain/model/companyresourcetype"
	"time4book/internal/app/core/domain/model/user"
)

type Facade struct {
	Create *Create
	List   *List
	Update *Update
	Delete *Delete
}

func NewFacade(repo companyresourcetype.Repo, userRepo user.UserRepo) *Facade {
	return &Facade{
		Create: &Create{repo: repo, userRepo: userRepo},
		List:   &List{repo: repo, userRepo: userRepo},
		Update: &Update{repo: repo, userRepo: userRepo},
		Delete: &Delete{repo: repo, userRepo: userRepo},
	}
}
