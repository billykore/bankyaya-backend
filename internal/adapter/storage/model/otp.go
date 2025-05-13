package model

import (
	"time"
)

type OTP struct {
	ID         int `gorm:"primaryKey"`
	Code       string
	UserID     int
	Purpose    string
	Channel    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	VerifiedAt time.Time
	ExpiredAt  time.Time

	User *User `gorm:"foreignKey:UserID"`
}
