package reservation

type ReservationStatus string

const (
	StatusActive           ReservationStatus = "active"
	StatusCancelled        ReservationStatus = "cancelled"
	StatusCancelledByAdmin ReservationStatus = "cancelled_by_admin"
	StatusCompleted        ReservationStatus = "completed"
)

func (s ReservationStatus) String() string { return string(s) }
