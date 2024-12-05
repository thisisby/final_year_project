package helpers

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	newValidator := validator.New()
	return &Validator{
		validator: newValidator,
	}
}

func (v *Validator) Validate(in any) error {
	return v.validator.Struct(in)
}
