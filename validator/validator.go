package validator

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var regexpMySQLErrorCode = regexp.MustCompile(`^Error ([0-9]{4})`)
var regexpMySQLErrorValue = regexp.MustCompile(`'(.*?)'`)

// ValidationError represents an error's specification.
type ValidationError struct {
	Text string `json:"text"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// NewValidationErrors create instance of ValidationError.
func NewValidationErrors(messages ...string) []ValidationError {
	var validationErrors []ValidationError

	for _, t := range messages {
		validationErrors = append(validationErrors, ValidationError{t})
	}

	return validationErrors
}

// Validate validates a structs fields.
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// FormattedValidationError returns formatted errors.
func FormattedValidationError(err error) []ValidationError {
	var validationErrors []ValidationError

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			t := fmt.Sprintf("%s must exist", e.Field())
			validationErrors = append(validationErrors, ValidationError{t})
		case "hexcolor":
			t := fmt.Sprintf("%s must be hexcolor", e.Field())
			validationErrors = append(validationErrors, ValidationError{t})
		case "max":
			t := fmt.Sprintf("%s is too long (maximum is %s characters)", e.Field(), e.Param())
			validationErrors = append(validationErrors, ValidationError{t})
		case "min":
			t := fmt.Sprintf("%s is too short (minimum is %s characters", e.Field(), e.Param())
			validationErrors = append(validationErrors, ValidationError{t})
		case "eqfield":
			t := fmt.Sprintf("%s must be equal to %s", e.Field(), e.Param())
			validationErrors = append(validationErrors, ValidationError{t})
		}
	}

	return validationErrors
}

// FormattedMySQLError returns formatted errors.
// found by MySQL error code.
func FormattedMySQLError(err error) []ValidationError {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
		}
	}()

	switch regexpMySQLErrorCode.FindStringSubmatch(err.Error())[1] {
	case "1062":
		v := strings.ReplaceAll(regexpMySQLErrorValue.FindString(err.Error()), "'", "")
		return NewValidationErrors(fmt.Sprintf("%s has already been taken", v))
	default:
		return nil
	}
}
