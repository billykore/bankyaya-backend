package dto

import (
	"time"

	"go.bankyaya.org/app/backend/internal/domain/user"
)

type LoginRequest struct {
	Phone      string `json:"phone" validate:"required,phonenumber"`
	Password   string `json:"password" validate:"required"`
	DeviceId   string `json:"deviceId" validate:"required"`
	FirebaseId string `json:"firebaseId" validate:"required"`
}

func (r *LoginRequest) ToUser() *user.User {
	return &user.User{
		CIF:           "",
		Password:      "",
		AccountNumber: "",
		FullName:      "",
		Email:         "",
		PhoneNumber:   r.Phone,
		NIK:           "",
		Device:        nil,
	}
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expiredTime"`
}

func NewLoginResponse(token *user.Token) *LoginResponse {
	return &LoginResponse{
		Token:     token.AccessToken,
		ExpiredAt: token.ExpiresAt,
	}
}
