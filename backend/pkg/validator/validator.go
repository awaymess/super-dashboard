package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init initializes the validator
func Init() {
	validate = validator.New()
}

// Get returns the validator instance
func Get() *validator.Validate {
	return validate
}

// Validate validates a struct
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}
