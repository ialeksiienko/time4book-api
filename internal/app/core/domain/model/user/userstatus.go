package user

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

func (s UserStatus) String() string { return string(s) }
