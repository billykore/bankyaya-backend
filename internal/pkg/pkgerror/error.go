package pkgerror

import (
	"go.bankyaya.org/app/backend/internal/pkg/codes"
)

// Error represents domain error.
type Error struct {
	// Code is the error code.
	Code codes.Code
	// Err is the error.
	Err error
}

// New returns new Error.
func New(c codes.Code, err error) error {
	return &Error{
		Code: c,
		Err:  err,
	}
}

func (err *Error) Error() string {
	return err.Err.Error()
}
