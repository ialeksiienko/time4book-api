package resource

import "errors"

var ErrInvalidType = errors.New("invalid resource type")
var ErrInvalidDetails = errors.New("invalid resource details")
var ErrInvalidServiceInterval = errors.New("invalid service interval")
var ErrInvalidServiceReason = errors.New("invalid service reason")
var ErrNotBookable = errors.New("resource is not bookable")
var ErrInvalidBookingInterval = errors.New("invalid booking interval")
var ErrReservationTooLong = errors.New("reservation exceeds maximum allowed time")
