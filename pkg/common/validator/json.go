package validator

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

type JSONValidator struct {
	validate *validator.Validate
}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{
		validate: validator.New(),
	}
}

// validates the given struct as with the rules defined by https://godoc.org/github.com/go-playground/validator
func (j JSONValidator) Validate(data interface{}) error {
	err := j.validate.Struct(data)
	if err != nil {
		return fmt.Errorf("%w", err.(validator.ValidationErrors))
	}
	return nil
}
