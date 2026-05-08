package resource

type ResourceType string

const (
	TypeRoom       ResourceType = "room"
	TypeCar        ResourceType = "car"
	TypeMultimedia ResourceType = "multimedia"
	TypeEquipment  ResourceType = "equipment"
)

func (t ResourceType) String() string {
	return string(t)
}
