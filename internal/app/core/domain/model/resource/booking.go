package resource

import "time"

func (r *Resource) IsBookable() bool { return r.status == StatusActive }

func (r *Resource) ValidateBookingInterval(start, end time.Time) error {
	if !r.IsBookable() {
		return ErrNotBookable
	}
	if start.IsZero() || !end.After(start) {
		return ErrInvalidBookingInterval
	}
	if r.maxReservationMinutes != nil {
		duration := end.Sub(start)
		minutes := duration / time.Minute
		if minutes > time.Duration(*r.maxReservationMinutes) ||
			(minutes == time.Duration(*r.maxReservationMinutes) && duration%time.Minute > 0) {
			return ErrReservationTooLong
		}
	}
	return nil
}
