package booking

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ListFilter struct {
	CompanyID  *uuid.UUID
	UserID     *uuid.UUID
	ResourceID *uuid.UUID
	Status     *BookingStatus
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type BookingRepo interface {
	Create(ctx context.Context, b *Booking) error
	ByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	List(ctx context.Context, f ListFilter) ([]*Booking, int64, error)
	ListByResourceID(ctx context.Context, resourceID uuid.UUID, from, to *time.Time, page, limit int) ([]*Booking, int64, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]*Booking, int64, error)
	Update(ctx context.Context, b *Booking) error
	ActiveByResourceIDInRange(ctx context.Context, resourceID uuid.UUID, from, to time.Time, excludeID *uuid.UUID) ([]*Booking, error)
}
