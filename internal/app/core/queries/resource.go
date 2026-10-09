package queries

import (
	"time"

	"github.com/google/uuid"
)

type CompanyResourceType struct {
	ID      uuid.UUID
	Name    string
	IconKey string
}

type Resource struct {
	ID                    uuid.UUID
	CompanyID             uuid.UUID
	Name                  string
	Type                  string
	Description           string
	Location              string
	MaxReservationMinutes *int
	AvailableFrom         *string
	AvailableTo           *string
	Status                string
	UnavailableFrom       *time.Time
	UnavailableTo         *time.Time
	UnavailableReason     *string
	CompanyResourceType   *CompanyResourceType
}
