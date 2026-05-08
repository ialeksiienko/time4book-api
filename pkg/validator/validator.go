package validator

import "github.com/go-playground/validator/v10"

type ValidationErrors validator.ValidationErrors

func (v ValidationErrors) Error() string {
	return validator.ValidationErrors(v).Error()
}

type Facade struct {
	*validator.Validate
}

func New() *Facade {
	return &Facade{validator.New(validator.WithRequiredStructEnabled())}
}
