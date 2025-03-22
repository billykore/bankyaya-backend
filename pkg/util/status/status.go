package status

import (
	"fmt"

	"go.bankyaya.org/app/backend/pkg/util/codes"
)

// Status represents domain error.
type Status struct {
	// Code is the error code.
	Code codes.Code
	// Err is the error.
	Err error
}

// Error returns new Status.
func Error(c codes.Code, err error) error {
	return &Status{
		Code: c,
		Err:  err,
	}
}

// Errorf returns new formatted Status.
func Errorf(c codes.Code, format string, a ...any) error {
	return Error(c, fmt.Errorf(format, a...))
}

func (s *Status) Error() string {
	return s.Err.Error()
}
