package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidPhoneNumber(t *testing.T) {
	v := validator.New()
	err := v.RegisterValidation("phonenumber", ValidPhoneNumber)
	assert.NoError(t, err)

	isErr := func(err error) bool {
		return err != nil
	}

	tests := []struct {
		name string
		args string
		want bool
	}{
		{name: "valid phone number", args: "+6281338442777", want: false},
		{name: "missing + in phone number", args: "6281338442777", want: true},
		{name: "+ sign in the end of phone number", args: "6281338442777+", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isErr(v.Var(tt.args, "phonenumber")), "ValidPhoneNumber(%v)", tt.args)
		})
	}
}
