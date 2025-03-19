package auth

import (
	"go.bankyaya.org/app/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	log *logger.Logger
}

func NewBcryptHasher(log *logger.Logger) *BcryptHasher {
	return &BcryptHasher{
		log: log,
	}
}

// Hash generates a bcrypt hash of the password.
func (b *BcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		b.log.Errorf("GenerateFromPassword error: %v", err)
		return "", err
	}
	return string(hashed), err
}

// Compare checks if the given password matches the stored hashed password.
func (b *BcryptHasher) Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		b.log.Errorf("CompareHashAndPassword error: %v", err)
		return false
	}
	return true
}
