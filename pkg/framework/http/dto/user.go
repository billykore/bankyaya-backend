package dto

import (
	"go.bankyaya.org/app/backend/pkg/entity"
)

type LoginRequest struct {
	Phone      string `json:"phone" validate:"required,phonenumber"`
	Password   string `json:"password" validate:"required"`
	DeviceId   string `json:"deviceId" validate:"required"`
	FirebaseId string `json:"firebaseId" validate:"required"`
}

func (r *LoginRequest) ToUser() *entity.User {
	return &entity.User{
		PhoneNumber: r.Phone,
		AuthData: entity.AuthData{
			Password:   r.Password,
			DeviceId:   r.DeviceId,
			FirebaseId: r.FirebaseId,
		},
	}
}

type LoginResponse struct {
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}

func NewLoginResponse(token *entity.Token) *LoginResponse {
	return &LoginResponse{
		Token:       token.AccessToken,
		ExpiredTime: token.ExpiredTime,
	}
}
