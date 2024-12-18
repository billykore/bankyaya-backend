package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationErr(t *testing.T) {
	type s struct {
		Required string `validate:"required"`
		Len      string `validate:"len=6"`
		Email    string `validate:"email"`
		Number   string `validate:"number"`
	}

	type test struct {
		name    string
		s       s
		wantMsg string
	}

	v := New()

	for _, tt := range []test{
		{
			name: "required validation error message",
			s: s{
				Required: "",
				Len:      "123456",
				Email:    "billy@kore.com",
				Number:   "7",
			},
			wantMsg: "Required is required",
		},
		{
			name: "len validation error message",
			s: s{
				Required: "required",
				Len:      "123",
				Email:    "billy@kore.com",
				Number:   "7",
			},
			wantMsg: "Len length must be 6",
		},
		{
			name: "email validation error message",
			s: s{
				Required: "required",
				Len:      "123456",
				Email:    "asdf",
				Number:   "7",
			},
			wantMsg: "Email is not a valid email",
		},
		{
			name: "number validation error message",
			s: s{
				Required: "required",
				Len:      "123456",
				Email:    "billy@kore.com",
				Number:   "asdf",
			},
			wantMsg: "Number must be a number",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.s)
			assert.Equal(t, tt.wantMsg, err.Error())
		})
	}
}
