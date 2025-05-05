package otp

import (
	"errors"
	"math/rand"
)

const digits = "0123456789"

type OTP struct {
}

func NewOTP() *OTP {
	return &OTP{}
}

func (o *OTP) Generate(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid OTP length")
	}
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp), nil
}
