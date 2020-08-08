package validator

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var regexpMySQLErrorCode = regexp.MustCompile(`^Error ([0-9]{4})`)
var regexpMySQLErrorValue = regexp.MustCompile(`'(.*?)'`)

// ValidationError represents an error's specification.
type ValidationError struct {
	Text string `json:"text"`
}

var (
	validate *validator.Validate
	locale   map[string]map[string]string
)

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		fieldName := f.Tag.Get("translationField")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	viper.SetConfigName("ja")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./validator/locales")
	viper.AddConfigPath("../validator/locales")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed read locale file: %v", err)
	}

	if err := viper.Unmarshal(&locale); err != nil {
		log.Printf("failed unmarshal config file: %v", err)
	}
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
		f, p := translateFieldError(e)

		switch e.Tag() {
		case "required":
			t := ErrorRequired(f)
			validationErrors = append(validationErrors, ValidationError{t})
		case "hexcolor":
			t := ErrorHexcolor(f)
			validationErrors = append(validationErrors, ValidationError{t})
		case "max":
			t := ErrorTooLong(f, p)
			validationErrors = append(validationErrors, ValidationError{t})
		case "min":
			t := ErrorTooShort(f, p)
			validationErrors = append(validationErrors, ValidationError{t})
		case "eqfield":
			t := ErrorEqualField(f, p)
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
		return NewValidationErrors(ErrorAlreadyBeenTaken)
	case "1452":
		return NewValidationErrors(ErrorForeignKeyConstraintFailed)
	default:
		return nil
	}
}

func translateFieldError(e validator.FieldError) (string, string) {
	ns := strings.Split(e.Namespace(), ".")
	// ns: expect [StrunctName, FieldName]

	if len(ns) < 2 {
		return e.Field(), e.Param()
	}

	s := strings.ToLower(ns[0])
	f := strings.ToLower(ns[1])
	p := strings.ToLower(e.Param())

	if s, ok1 := locale[s]; ok1 {
		if f, ok2 := s[f]; ok2 {
			if p, ok3 := s[p]; ok3 {
				return f, p
			}
			return f, e.Param()
		}
	}
	return e.Field(), e.Param()
}
