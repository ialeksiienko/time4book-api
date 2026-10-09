package reservationport

import (
	"context"
	"errors"
	"time"
	domain "time4book/internal/app/core/domain/model/reservation"

	"github.com/google/uuid"
)

var ErrTimeSlotBooked = errors.New("time slot is already booked")

type Filter struct {
	CompanyID  *uuid.UUID
	UserID     *uuid.UUID
	ResourceID *uuid.UUID
	Status     *domain.ReservationStatus
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type Repository interface {
	Create(ctx context.Context, r *domain.Reservation) error
	Retrieve(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	Get(ctx context.Context, f Filter) ([]*domain.Reservation, int64, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID, companyID uuid.UUID, from, to *time.Time, page, limit int) ([]*domain.Reservation, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, page, limit int) ([]*domain.Reservation, int64, error)
	Update(ctx context.Context, r *domain.Reservation) error
}
