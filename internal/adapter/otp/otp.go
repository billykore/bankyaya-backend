package otp

import (
	"crypto/rand"
	"errors"
	"math/big"
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
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}
