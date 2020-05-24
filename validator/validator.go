package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Text string `json:"text"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func MakeErrors(messages ...string) []ValidationError {
	var validationErrors []ValidationError

	for _, t := range messages {
		validationErrors = append(validationErrors, ValidationError{t})
	}

	return validationErrors
}

func Validate(s interface{}) error {
	return validate.Struct(s)
}

func ValidationMessages(err error) []ValidationError {
	var validationErrors []ValidationError

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			t := fmt.Sprintf("%s must exist", e.Field())
			validationErrors = append(validationErrors, ValidationError{t})
		}
	}

	return validationErrors
}
