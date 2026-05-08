package kernel

import (
	"errors"
	"time"
)

type TimeRange struct {
	fromDate time.Time
	toDate   time.Time
}

func NewTimeRange(
	fromDate time.Time,
	toDate time.Time,
) (TimeRange, error) {
	if fromDate.IsZero() || toDate.IsZero() {
		return TimeRange{}, errors.New("dates cannot be zero")
	}

	if toDate.Before(fromDate) || toDate.Equal(fromDate) {
		return TimeRange{}, errors.New("to_date must be after from_date")
	}

	return TimeRange{
		fromDate: fromDate,
		toDate:   toDate,
	}, nil
}

func (t TimeRange) Overlaps(other TimeRange) bool {
	return t.fromDate.Before(other.toDate) && t.toDate.After(other.fromDate)
}
