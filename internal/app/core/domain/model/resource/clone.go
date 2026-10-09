package resource

func clone[T any](value *T) *T {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
