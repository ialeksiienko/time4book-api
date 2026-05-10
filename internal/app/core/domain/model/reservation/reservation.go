package reservation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	id          uuid.UUID
	userID      uuid.UUID
	companyID   uuid.UUID
	resourceID  uuid.UUID
	description *string
	startDate   time.Time
	endDate     time.Time
	status      ReservationStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewReservation(
	userID uuid.UUID,
	companyID uuid.UUID,
	resourceID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
	description *string,
) (*Reservation, error) {
	if !endDate.After(startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	return &Reservation{
		id:          uuid.New(),
		userID:      userID,
		companyID:   companyID,
		resourceID:  resourceID,
		description: description,
		startDate:   startDate,
		endDate:     endDate,
		status:      StatusActive,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}, nil
}

func (r *Reservation) Cancel() error {
	if r.status == StatusCompleted {
		return errors.New("cannot cancel a completed reservation")
	}
	r.status = StatusCancelled
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *Reservation) CancelByAdmin() error {
	if r.status == StatusCompleted {
		return errors.New("cannot cancel a completed reservation")
	}
	r.status = StatusCancelledByAdmin
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *Reservation) Complete() error {
	if r.status != StatusActive {
		return errors.New("reservation can only be completed from active status")
	}
	r.status = StatusCompleted
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *Reservation) ID() uuid.UUID             { return r.id }
func (r *Reservation) UserID() uuid.UUID         { return r.userID }
func (r *Reservation) CompanyID() uuid.UUID      { return r.companyID }
func (r *Reservation) ResourceID() uuid.UUID     { return r.resourceID }
func (r *Reservation) Description() *string      { return r.description }
func (r *Reservation) StartDate() time.Time      { return r.startDate }
func (r *Reservation) EndDate() time.Time        { return r.endDate }
func (r *Reservation) Status() ReservationStatus { return r.status }
func (r *Reservation) CreatedAt() time.Time      { return r.createdAt }
func (r *Reservation) UpdatedAt() time.Time      { return r.updatedAt }
func (r *Reservation) CancelledAt() *time.Time {
	if r.status == StatusCancelled || r.status == StatusCancelledByAdmin {
		return &r.updatedAt
	}
	return nil
}

type Props struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CompanyID   uuid.UUID
	ResourceID  uuid.UUID
	Description *string
	StartDate   time.Time
	EndDate     time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Reconstitute(p *Props) *Reservation {
	return &Reservation{
		id:          p.ID,
		userID:      p.UserID,
		companyID:   p.CompanyID,
		resourceID:  p.ResourceID,
		description: p.Description,
		startDate:   p.StartDate,
		endDate:     p.EndDate,
		status:      ReservationStatus(p.Status),
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}
