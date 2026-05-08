package resource

type ResourceStatus string

const (
	StatusActive    ResourceStatus = "active"
	StatusInactive  ResourceStatus = "inactive"
	StatusInService ResourceStatus = "in_service"
)

func (s ResourceStatus) String() string { return string(s) }
