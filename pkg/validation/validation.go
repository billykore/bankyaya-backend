package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a struct that provides validation functionality for request data.
type Validator struct {
	v *validator.Validate
}

// New initializes and returns a new instance of Validator.
func New() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

// Validate checks the validity of the input data based on struct tags and returns an error if any.
func (v *Validator) Validate(req any) error {
	err := v.v.Struct(req)
	return validationErr(err)
}

// validationErr processes validation errors and returns a formatted error message.
func validationErr(err error) error {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	return joinValidationErrors(ve)
}

// joinValidationErrors converts a slice of FieldError into a concatenated error.
func joinValidationErrors(validationErrors validator.ValidationErrors) error {
	var errList []string
	for _, fieldError := range validationErrors {
		errList = append(errList, errMessage(fieldError))
	}
	return fmt.Errorf(strings.Join(errList, "; "))
}

// tagMessages maps validation tags to corresponding error message templates.
var tagMessages = map[string]string{
	"required": "%s is required",
	"email":    "%s is not a valid email",
	"len":      "%s length must be %s",
	"min":      "%s minimum length must be %s",
	"number":   "%s must be a number",
}

// errMessage formats and returns error messages based on the field validation error type.
func errMessage(fe validator.FieldError) string {
	if messageTemplate, exists := tagMessages[fe.Tag()]; exists {
		return formatMessage(messageTemplate, fe.Field(), fe.Param())
	}
	return formatMessage("%s format is invalid", fe.Field())
}

// formatMessage is a helper function to format messages with parameters.
func formatMessage(template, field string, param ...string) string {
	if len(param) > 0 && param[0] != "" {
		return fmt.Sprintf(template, field, param[0])
	}
	return fmt.Sprintf(template, field)
}
