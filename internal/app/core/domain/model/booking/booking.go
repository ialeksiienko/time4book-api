package booking

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	id          uuid.UUID
	userID      uuid.UUID
	resourceID  uuid.UUID
	description *string
	startDate   time.Time
	endDate     time.Time
	status      BookingStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBooking(
	userID uuid.UUID,
	resourceID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
	description *string,
) (*Booking, error) {
	if !endDate.After(startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	return &Booking{
		id:          uuid.New(),
		userID:      userID,
		resourceID:  resourceID,
		description: description,
		startDate:   startDate,
		endDate:     endDate,
		status:      StatusActive,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}, nil
}

func (b *Booking) Cancel() error {
	if b.status == StatusCompleted {
		return errors.New("cannot cancel a completed booking")
	}
	b.status = StatusCancelled
	b.updatedAt = time.Now().UTC()
	return nil
}

func (b *Booking) CancelByAdmin() error {
	if b.status == StatusCompleted {
		return errors.New("cannot cancel a completed booking")
	}
	b.status = StatusCancelledByAdmin
	b.updatedAt = time.Now().UTC()
	return nil
}

func (b *Booking) Complete() error {
	if b.status != StatusActive {
		return errors.New("booking can only be completed from active status")
	}
	b.status = StatusCompleted
	b.updatedAt = time.Now().UTC()
	return nil
}

func (b *Booking) ID() uuid.UUID         { return b.id }
func (b *Booking) UserID() uuid.UUID     { return b.userID }
func (b *Booking) ResourceID() uuid.UUID { return b.resourceID }
func (b *Booking) Description() *string  { return b.description }
func (b *Booking) StartDate() time.Time  { return b.startDate }
func (b *Booking) EndDate() time.Time    { return b.endDate }
func (b *Booking) Status() BookingStatus { return b.status }
func (b *Booking) CreatedAt() time.Time  { return b.createdAt }
func (b *Booking) UpdatedAt() time.Time  { return b.updatedAt }
func (b *Booking) CancelledAt() *time.Time {
	if b.status == StatusCancelled || b.status == StatusCancelledByAdmin {
		return &b.updatedAt
	}
	return nil
}

type Props struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ResourceID  uuid.UUID
	Description *string
	StartDate   time.Time
	EndDate     time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Reconstitute(p *Props) *Booking {
	return &Booking{
		id:          p.ID,
		userID:      p.UserID,
		resourceID:  p.ResourceID,
		description: p.Description,
		startDate:   p.StartDate,
		endDate:     p.EndDate,
		status:      BookingStatus(p.Status),
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}
