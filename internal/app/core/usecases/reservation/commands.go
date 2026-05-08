package reservationcommands

import (
	"log/slog"
	"time4book/internal/app/core/domain/model/booking"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	"time4book/pkg/validator"
)

type Facade struct {
	Create         *Create
	List           *List
	ListByResource *ListByResource
	ListMy         *ListMy
	Cancel         *Cancel
}

func NewFacade(
	urepo user.UserRepo,
	resrepo resource.ResourceRepo,
	brepo booking.BookingRepo,
	v *validator.Facade,
	log *slog.Logger,
) *Facade {
	return &Facade{
		Create:         newCreate(urepo, resrepo, brepo, v, log.With(slog.String("reservation", "create"))),
		List:           newList(brepo, log.With(slog.String("reservation", "list"))),
		ListByResource: newListByResource(brepo, log.With(slog.String("reservation", "listByResource"))),
		ListMy:         newListMy(brepo, log.With(slog.String("reservation", "listMy"))),
		Cancel:         newCancel(urepo, brepo, log.With(slog.String("reservation", "cancel"))),
	}
}
