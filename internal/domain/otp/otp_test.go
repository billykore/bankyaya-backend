package otp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPurpose(t *testing.T) {
	purpose := "login"
	assert.IsType(t, Purpose("login"), NewPurpose(purpose))
}

func TestPurposeMessage(t *testing.T) {
	purpose := PurposeRegister
	assert.Equal(t, "This is your OTP for register", purpose.Message())
}

func TestPurposeString(t *testing.T) {
	purpose := PurposeRegister
	assert.Equal(t, "register", purpose.String())
}

func TestNewChannel(t *testing.T) {
	channel := "sms"
	assert.IsType(t, Channel("sms"), NewChannel(channel))
}

func TestChannelString(t *testing.T) {
	channel := ChannelSMS
	assert.Equal(t, "sms", channel.String())
}

func TestOTPEqual(t *testing.T) {
	otp1 := &OTP{
		ID:      1,
		Code:    "123456",
		Purpose: PurposeLogin,
		Channel: ChannelEmail,
		User: &User{
			ID:    1,
			Name:  "John Doe",
			Email: "jd@email.com",
			Phone: "123",
		},
		CreatedAt:  time.Time{},
		ExpiredAt:  time.Time{},
		VerifiedAt: time.Time{},
	}
	otp2 := &OTP{
		ID:      2,
		Code:    "123456",
		Purpose: PurposeLogin,
		Channel: ChannelEmail,
		User: &User{
			ID:    1,
			Name:  "John Doe",
			Email: "jd@email.com",
			Phone: "123",
		},
		CreatedAt:  time.Time{},
		ExpiredAt:  time.Time{},
		VerifiedAt: time.Time{},
	}

	assert.False(t, otp1.Equal(otp2))
	assert.False(t, otp2.Equal(nil))
}

func TestNewUser(t *testing.T) {
	user := NewUser(1, "John Doe", "jd@email.com", "123")
	assert.Equal(t, &User{
		ID:    1,
		Name:  "John Doe",
		Email: "jd@email.com",
		Phone: "123",
	}, user)
}

func TestUserEqual(t *testing.T) {
	user1 := NewUser(1, "John Doe", "jd@email.com", "123")
	user2 := NewUser(2, "John Doe", "jd@email.com", "123")
	assert.False(t, user1.Equal(user2))
	assert.False(t, user2.Equal(nil))
}
