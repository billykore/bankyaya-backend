// Package otp provides structures and functionality for managing One-Time Passwords (OTP)
// used in authentication and verification processes.
package otp

import (
	"time"
)

// Purpose defines the purpose of an action, such as login or registration.
type Purpose string

const (
	PurposeLogin    Purpose = "login"
	PurposeRegister Purpose = "register"
)

// NewPurpose creates a new Purpose from the given string.
func NewPurpose(purpose string) Purpose {
	return Purpose(purpose)
}

// Message returns a formatted OTP message based on the Purpose value.
func (p Purpose) Message() string {
	return "This is your OTP for " + string(p)
}

// String converts the Purpose value to its string representation.
func (p Purpose) String() string {
	return string(p)
}

// Channel represents a type of communication medium, such as email or SMS.
type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
)

// NewChannel creates and returns a new Channel instance from the provided string.
func NewChannel(channel string) Channel {
	return Channel(channel)
}

func (c Channel) String() string {
	return string(c)
}

// OTP represents a one-time password used for authentication or verification.
type OTP struct {
	ID         uint64
	Code       string
	Purpose    Purpose
	Channel    Channel
	User       *User
	CreatedAt  time.Time
	ExpiredAt  time.Time
	VerifiedAt time.Time
}

// IsExpired returns true if the OTP is expired.
func (o *OTP) IsExpired(now time.Time) bool {
	return now.After(o.ExpiredAt)
}

// IsVerified returns true if the OTP has already been verified.
func (o *OTP) IsVerified() bool {
	return !o.VerifiedAt.IsZero()
}

// Equal compares two OTP instances and returns true if their ID,
// Code, Purpose, Channel, and UserID match.
func (o *OTP) Equal(other *OTP) bool {
	if o != nil && other != nil {
		return o.ID == other.ID &&
			o.Code == other.Code &&
			o.Purpose == other.Purpose &&
			o.Channel == other.Channel &&
			o.User.Equal(other.User)
	}
	return false
}

// User represents a user with an ID, name, email, and phone.
type User struct {
	ID    int
	Name  string
	Email string
	Phone string
}

// NewUser creates and returns a new User instance with the provided ID, name, email, and phone.
func NewUser(id int, name, email, phone string) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Phone: phone,
	}
}

// Equal compares the current User with another User and returns true if all fields are equal.
func (u *User) Equal(other *User) bool {
	if u != nil && other != nil {
		return u.ID == other.ID &&
			u.Name == other.Name &&
			u.Email == other.Email &&
			u.Phone == other.Phone
	}
	return false
}
