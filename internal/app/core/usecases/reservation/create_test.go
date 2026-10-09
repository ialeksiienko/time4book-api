package reservationcommands

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"time4book/internal/app/core/domain/model/reservation"
	"time4book/internal/app/core/domain/model/resource"
	"time4book/internal/app/core/domain/model/user"
	port "time4book/internal/app/core/ports/in/reservation"
	reservationout "time4book/internal/app/core/ports/out/reservation"
	resourceout "time4book/internal/app/core/ports/out/resource"
	userout "time4book/internal/app/core/ports/out/user"
	"time4book/pkg/validator"

	"github.com/google/uuid"
)

type createUserRepo struct {
	userout.Repository
	account *user.User
}

func (r createUserRepo) Retrieve(context.Context, uuid.UUID) (*user.User, error) {
	return r.account, nil
}

type createResourceRepo struct {
	resourceout.Repository
	resource *resource.Resource
}

func (r createResourceRepo) Retrieve(context.Context, uuid.UUID) (*resource.Resource, error) {
	return r.resource, nil
}

type createReservationRepo struct {
	reservationout.Repository
	created bool
}

func (r *createReservationRepo) Create(context.Context, *reservation.Reservation) error {
	r.created = true
	return nil
}

func TestCreateRejectsCompanyMismatchAndInactiveUser(t *testing.T) {
	companyID, otherCompanyID := uuid.New(), uuid.New()
	start := time.Date(2026, time.October, 9, 10, 0, 0, 0, time.UTC)

	for _, tt := range []struct {
		name            string
		userCompany     uuid.UUID
		resourceCompany uuid.UUID
		requestCompany  uuid.UUID
		userStatus      string
	}{
		{"request company differs from user", companyID, companyID, otherCompanyID, "active"},
		{"request company differs from resource", companyID, otherCompanyID, companyID, "active"},
		{"inactive user", companyID, companyID, companyID, "inactive"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			accountID, resourceID := uuid.New(), uuid.New()
			account := user.Reconstitute(&user.Props{ID: accountID, CompanyID: tt.userCompany, Status: tt.userStatus})
			bookableResource := resource.Reconstitute(&resource.Props{ID: resourceID, CompanyID: tt.resourceCompany, Status: "active"})
			bookings := &createReservationRepo{}
			command := newCreate(createUserRepo{account: account}, createResourceRepo{resource: bookableResource}, bookings, validator.New(), slog.New(slog.NewTextHandler(io.Discard, nil)))
			_, err := command.Execute(context.Background(), &port.CreateRequest{
				InitiatorID: accountID,
				CompanyID:   tt.requestCompany,
				ResourceID:  resourceID,
				StartDate:   start,
				EndDate:     start.Add(time.Hour),
			})
			if !errors.Is(err, user.ErrUnauthorized) {
				t.Fatalf("got error %v, want unauthorized", err)
			}
			if bookings.created {
				t.Fatal("created reservation without company access")
			}
		})
	}
}
