package status

import (
	"fmt"

	"go.bankyaya.org/app/backend/pkg/codes"
)

// Status represents domain error.
type Status struct {
	// Code is the error code.
	Code codes.Code
	// Message is the error message.
	Message string
}

// Error returns new Status.
func Error(c codes.Code, msg string) error {
	return &Status{
		Code:    c,
		Message: msg,
	}
}

// Errorf returns new formatted Status.
func Errorf(c codes.Code, format string, a ...any) error {
	return Error(c, fmt.Sprintf(format, a...))
}

func (s *Status) Error() string {
	return s.Message
}
