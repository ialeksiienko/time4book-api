package resource

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewResourceRejectsInvalidDetails(t *testing.T) {
	companyID := uuid.New()
	customTypeID := uuid.New()
	zero := 0
	negative := -1
	from := "09:00"
	to := "17:00"
	invalidHour := "9:00"
	early := "08:00"
	cases := []struct {
		name    string
		details Details
		wantErr error
	}{
		{"blank name", Details{Name: "  ", Type: TypeRoom}, ErrInvalidDetails},
		{"unknown type", Details{Name: "Room", Type: "unknown"}, ErrInvalidType},
		{"custom without id", Details{Name: "Room", Type: TypeCustom}, ErrInvalidType},
		{"custom with empty id", Details{Name: "Room", Type: TypeCustom, CompanyResourceTypeID: &uuid.Nil}, ErrInvalidType},
		{"builtin with custom id", Details{Name: "Room", Type: TypeRoom, CompanyResourceTypeID: &customTypeID}, ErrInvalidType},
		{"zero duration", Details{Name: "Room", Type: TypeRoom, MaxReservationMinutes: &zero}, ErrInvalidDetails},
		{"negative duration", Details{Name: "Room", Type: TypeRoom, MaxReservationMinutes: &negative}, ErrInvalidDetails},
		{"one daily hour", Details{Name: "Room", Type: TypeRoom, AvailableFrom: &from}, ErrInvalidDetails},
		{"invalid daily hour", Details{Name: "Room", Type: TypeRoom, AvailableFrom: &invalidHour, AvailableTo: &to}, ErrInvalidDetails},
		{"reversed daily hours", Details{Name: "Room", Type: TypeRoom, AvailableFrom: &from, AvailableTo: &early}, ErrInvalidDetails},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewResource(companyID, tc.details)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got error %v; want %v", err, tc.wantErr)
			}
		})
	}
}

func TestUpdateDetailsPreservesIdentityAndServiceState(t *testing.T) {
	res, err := NewResource(uuid.New(), Details{Name: "Old", Type: TypeRoom})
	if err != nil {
		t.Fatal(err)
	}
	from := time.Now().UTC().Add(time.Hour)
	to := from.Add(time.Hour)
	serviced, err := res.MarkInService("repair", from, &to)
	if err != nil {
		t.Fatal(err)
	}
	maxMinutes := 45
	updated, err := serviced.UpdateDetails(Details{Name: "New", Type: TypeCar, MaxReservationMinutes: &maxMinutes})
	if err != nil {
		t.Fatal(err)
	}
	if updated.ID() != res.ID() || updated.CompanyID() != res.CompanyID() || updated.CreatedAt() != res.CreatedAt() {
		t.Fatal("update changed identity or creation time")
	}
	if updated.Status() != StatusInService || !updated.UnavailableFrom().Equal(from) || !updated.UnavailableTo().Equal(to) || *updated.UnavailableReason() != "repair" {
		t.Fatal("update changed service state")
	}
	if updated.Name() != "New" || updated.Type() != TypeCar || *updated.MaxReservationMinutes() != maxMinutes {
		t.Fatal("update did not apply details")
	}
	if res.Name() != "Old" || serviced.Name() != "Old" || res.Status() != StatusActive {
		t.Fatal("update mutated a prior resource")
	}
	if updated.UpdatedAt().Before(serviced.UpdatedAt()) {
		t.Fatal("update time moved backwards")
	}
	_, err = serviced.UpdateDetails(Details{Name: " ", Type: TypeRoom})
	if !errors.Is(err, ErrInvalidDetails) || serviced.Name() != "Old" {
		t.Fatal("invalid update changed the resource")
	}
}

func TestServiceIntervalAndBookingDuration(t *testing.T) {
	maxMinutes := 30
	res, err := NewResource(uuid.New(), Details{Name: "Room", Type: TypeRoom, MaxReservationMinutes: &maxMinutes})
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now().UTC().Add(time.Hour)
	if err := res.ValidateBookingInterval(start, start.Add(30*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := res.ValidateBookingInterval(start, start.Add(30*time.Minute+time.Nanosecond)); !errors.Is(err, ErrReservationTooLong) {
		t.Fatalf("got error %v; want reservation too long", err)
	}
	if err := res.ValidateBookingInterval(start, start); !errors.Is(err, ErrInvalidBookingInterval) {
		t.Fatalf("got error %v; want invalid booking interval", err)
	}
	if _, err := res.MarkInService("repair", start, &start); !errors.Is(err, ErrInvalidServiceInterval) {
		t.Fatalf("got error %v; want invalid service interval", err)
	}
	if _, err := res.MarkInService(" ", start, nil); !errors.Is(err, ErrInvalidServiceReason) {
		t.Fatalf("got error %v; want invalid service reason", err)
	}
	serviced, err := res.MarkInService("repair", start, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := serviced.ValidateBookingInterval(start, start.Add(time.Minute)); !errors.Is(err, ErrNotBookable) {
		t.Fatalf("got error %v; want not bookable", err)
	}
	if res.Status() != StatusActive {
		t.Fatal("service mutated the original resource")
	}
	restored := serviced.Restore()
	if restored.Status() != StatusActive || restored.UnavailableFrom() != nil || restored.UnavailableReason() != nil {
		t.Fatal("restore did not clear service state")
	}
	if serviced.Status() != StatusInService {
		t.Fatal("restore mutated the original resource")
	}
}

func TestResourceCopiesInputAndOutputPointers(t *testing.T) {
	companyTypeID := uuid.New()
	originalTypeID := companyTypeID
	maxMinutes := 30
	availableFrom := "09:00"
	availableTo := "17:00"
	details := Details{
		Name: "Room", Type: TypeCustom, CompanyResourceTypeID: &companyTypeID,
		MaxReservationMinutes: &maxMinutes, AvailableFrom: &availableFrom, AvailableTo: &availableTo,
	}
	res, err := NewResource(uuid.New(), details)
	if err != nil {
		t.Fatal(err)
	}
	companyTypeID = uuid.New()
	maxMinutes = -1
	availableFrom = "invalid"
	availableTo = "invalid"
	if *res.CompanyResourceTypeID() != originalTypeID || *res.MaxReservationMinutes() != 30 ||
		*res.AvailableFrom() != "09:00" || *res.AvailableTo() != "17:00" {
		t.Fatal("constructor kept pointers owned by the caller")
	}
	*res.CompanyResourceTypeID() = uuid.Nil
	*res.MaxReservationMinutes() = -1
	*res.AvailableFrom() = "invalid"
	*res.AvailableTo() = "invalid"
	if *res.CompanyResourceTypeID() != originalTypeID || *res.MaxReservationMinutes() != 30 ||
		*res.AvailableFrom() != "09:00" || *res.AvailableTo() != "17:00" {
		t.Fatal("getter exposed resource state")
	}

	updatedMax := 45
	updatedFrom := "10:00"
	updatedTo := "18:00"
	updated, err := res.UpdateDetails(Details{
		Name: "Updated", Type: TypeRoom, MaxReservationMinutes: &updatedMax,
		AvailableFrom: &updatedFrom, AvailableTo: &updatedTo,
	})
	if err != nil {
		t.Fatal(err)
	}
	updatedMax = 0
	updatedFrom = "invalid"
	updatedTo = "invalid"
	if *updated.MaxReservationMinutes() != 45 || *updated.AvailableFrom() != "10:00" || *updated.AvailableTo() != "18:00" {
		t.Fatal("update kept pointers owned by the caller")
	}

	serviceFrom := time.Now().UTC().Add(time.Hour)
	serviceTo := serviceFrom.Add(time.Hour)
	serviced, err := updated.MarkInService("repair", serviceFrom, &serviceTo)
	if err != nil {
		t.Fatal(err)
	}
	expectedTo := serviceTo
	serviceTo = serviceTo.Add(time.Hour)
	*serviced.UnavailableFrom() = serviceFrom.Add(time.Hour)
	*serviced.UnavailableTo() = serviceTo
	*serviced.UnavailableReason() = "changed"
	if !serviced.UnavailableFrom().Equal(serviceFrom) || !serviced.UnavailableTo().Equal(expectedTo) ||
		*serviced.UnavailableReason() != "repair" {
		t.Fatal("service state shared input or output pointers")
	}
}

func TestReconstituteCopiesPointers(t *testing.T) {
	maxMinutes := 30
	serviceFrom := time.Now().UTC()
	serviceTo := serviceFrom.Add(time.Hour)
	reason := "repair"
	props := &Props{
		ID: uuid.New(), CompanyID: uuid.New(), Name: "Room", ResourceType: TypeRoom.String(),
		MaxReservationMinutes: &maxMinutes, Status: StatusInService.String(),
		UnavailableFrom: &serviceFrom, UnavailableTo: &serviceTo, UnavailableReason: &reason,
	}
	res := Reconstitute(props)
	maxMinutes = 0
	serviceFrom = serviceFrom.Add(time.Hour)
	serviceTo = serviceTo.Add(time.Hour)
	reason = "changed"
	if *res.MaxReservationMinutes() != 30 || res.UnavailableFrom().Equal(serviceFrom) ||
		res.UnavailableTo().Equal(serviceTo) || *res.UnavailableReason() != "repair" {
		t.Fatal("reconstituted resource kept pointers owned by the caller")
	}
}
