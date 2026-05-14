package booking

type BookingStatus string

const (
	StatusActive           BookingStatus = "active"
	StatusCancelled        BookingStatus = "cancelled"
	StatusCancelledByAdmin BookingStatus = "cancelled_by_admin"
	StatusCompleted        BookingStatus = "completed"
)

func (s BookingStatus) String() string { return string(s) }
