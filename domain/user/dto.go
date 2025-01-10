package user

type LoginRequest struct {
	Phone      string `json:"phone" validate:"required,phonenumber"`
	Password   string `json:"password" validate:"required"`
	DeviceId   string `json:"deviceId" validate:"required"`
	FirebaseId string `json:"firebaseId" validate:"required"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}
