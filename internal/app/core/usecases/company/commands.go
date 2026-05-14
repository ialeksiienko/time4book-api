package companycommands

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/booking"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"
	"time4book/pkg/validator"
)

type Facade struct {
	Create  *Create
	List    *List
	GetByID *GetByID
	Update  *Update
	Block   *Block
	Unblock *Unblock
	Delete  *Delete
}

func NewFacade(
	urepo user.UserRepo,
	crepo company.CompanyRepo,
	brepo booking.BookingRepo,
	tx ports.TxManager,
	v *validator.Facade,
	log *slog.Logger,
) *Facade {
	return &Facade{
		Create:  newCreate(crepo, tx, v, log.With(slog.String("company", "create"))),
		List:    newList(crepo, log.With(slog.String("company", "list"))),
		GetByID: newGetByID(crepo, log.With(slog.String("company", "getByID"))),
		Update:  newUpdate(urepo, crepo, v, log.With(slog.String("company", "update"))),
		Block:   newBlock(urepo, crepo, log.With(slog.String("company", "block"))),
		Unblock: newUnblock(urepo, crepo, log.With(slog.String("company", "unblock"))),
		Delete:  newDelete(urepo, crepo, brepo, tx, log.With(slog.String("company", "delete"))),
	}
}
