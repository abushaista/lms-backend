package validator

import (
	"github.com/go-playground/validator/v10"
)

type EchoValidator struct {
	v *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{v: validator.New()}
}

func (e *EchoValidator) Validate(i interface{}) error {
	return e.v.Struct(i)
}
