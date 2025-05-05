package model

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	Code       string
	UserID     int
	Purpose    string
	Channel    string
	VerifiedAt time.Time
	ExpiredAt  time.Time

	User *User `gorm:"foreignKey:UserID"`
}
