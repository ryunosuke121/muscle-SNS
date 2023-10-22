package echoValidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator is the validator for Echo framework.
type Validator struct {
	Validator *validator.Validate
}

// New returns a new validator for Echo framework.
func New(validator *validator.Validate) echo.Validator {
	return &Validator{
		Validator: validator,
	}
}

// Validate validates request body.
func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
