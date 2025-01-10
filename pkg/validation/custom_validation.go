package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidPhoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	exp := regexp.MustCompile(`^[+][0-9]+`)
	return exp.MatchString(phone)
}
