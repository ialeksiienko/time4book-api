package resource

import (
	"strings"
	"time"
)

func (r *Resource) MarkInService(reason string, from time.Time, to *time.Time) (*Resource, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, ErrInvalidServiceReason
	}
	if from.IsZero() || (to != nil && !to.After(from)) {
		return nil, ErrInvalidServiceInterval
	}
	updated := *r
	updated.status = StatusInService
	updated.unavailableFrom = &from
	updated.unavailableTo = clone(to)
	updated.unavailableReason = &reason
	updated.updatedAt = time.Now().UTC()
	return &updated, nil
}

func (r *Resource) Restore() *Resource {
	updated := *r
	updated.status = StatusActive
	updated.unavailableFrom = nil
	updated.unavailableTo = nil
	updated.unavailableReason = nil
	updated.updatedAt = time.Now().UTC()
	return &updated
}

func (r *Resource) Deactivate() *Resource {
	updated := *r
	updated.status = StatusInactive
	updated.updatedAt = time.Now().UTC()
	return &updated
}
