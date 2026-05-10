package reservation

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ListFilter struct {
	CompanyID  *uuid.UUID
	UserID     *uuid.UUID
	ResourceID *uuid.UUID
	Status     *ReservationStatus
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type ReservationRepo interface {
	Create(ctx context.Context, r *Reservation) error
	ByID(ctx context.Context, id uuid.UUID) (*Reservation, error)
	List(ctx context.Context, f ListFilter) ([]*Reservation, int64, error)
	ListByResourceID(ctx context.Context, resourceID uuid.UUID, from, to *time.Time, page, limit int) ([]*Reservation, int64, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Reservation, int64, error)
	Update(ctx context.Context, r *Reservation) error
	ActiveByResourceIDInRange(ctx context.Context, resourceID uuid.UUID, from, to time.Time, excludeID *uuid.UUID) ([]*Reservation, error)
}
