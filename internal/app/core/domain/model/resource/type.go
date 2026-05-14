package resource

type ResourceType string

const (
	TypeRoom       ResourceType = "room"
	TypeCar        ResourceType = "car"
	TypeMultimedia ResourceType = "multimedia"
	TypeEquipment  ResourceType = "equipment"
	// TypeCustom means display/booking semantics come from company_resource_types (per-company definition).
	TypeCustom ResourceType = "custom"
)

func IsBuiltInType(t ResourceType) bool {
	switch t {
	case TypeRoom, TypeCar, TypeMultimedia, TypeEquipment:
		return true
	default:
		return false
	}
}

func (t ResourceType) String() string {
	return string(t)
}
