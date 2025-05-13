package dto

import (
	"time"

	"go.bankyaya.org/app/backend/internal/domain/otp"
)

type OTPRequest struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Channel string `json:"channel"`
	Purpose string `json:"purpose"`
}

type OTPResponse struct {
	ID        int          `json:"id"`
	Code      string       `json:"code"`
	Channel   string       `json:"channel"`
	Purpose   string       `json:"purpose"`
	Recipient OTPRecipient `json:"recipient"`
	CreatedAt time.Time    `json:"createdAt"`
	ExpiredAt time.Time    `json:"expiredAt"`
}

type OTPRecipient struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func NewOTPResponse(otp *otp.OTP) *OTPResponse {
	return &OTPResponse{
		ID:      otp.ID,
		Code:    otp.Code,
		Channel: otp.Channel.String(),
		Purpose: otp.Purpose.String(),
		Recipient: OTPRecipient{
			Email: otp.User.Email,
			Phone: otp.User.Phone,
		},
		CreatedAt: otp.CreatedAt,
		ExpiredAt: otp.ExpiredAt,
	}
}

type VerifyOTPRequest struct {
	Code    string `json:"code"`
	UserID  int    `json:"userId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Channel string `json:"channel"`
	Purpose string `json:"purpose"`
}

func (r *VerifyOTPRequest) ToOTP() *otp.OTP {
	return &otp.OTP{
		Code:    r.Code,
		User:    otp.NewUser(r.UserID, r.Name, r.Email, r.Phone),
		Channel: otp.NewChannel(r.Channel),
		Purpose: otp.NewPurpose(r.Purpose),
	}
}
