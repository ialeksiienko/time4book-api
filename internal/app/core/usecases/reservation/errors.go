package reservationcommands

import "errors"

// ErrSlotAlreadyTaken is returned when the requested interval overlaps an active booking.
var ErrSlotAlreadyTaken = errors.New("wybrany termin jest już zajęty")

// ErrReservationStartTooFarInPast is returned when the booking starts earlier than allowed.
var ErrReservationStartTooFarInPast = errors.New("reservation start is too far in the past")
