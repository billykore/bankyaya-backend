package auth

import (
	"time"

	"go.bankyaya.org/app/backend/internal/core/user"
	"go.bankyaya.org/app/backend/pkg/data"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/security/token"
)

type JWT struct {
	log *logger.Logger
}

func NewJWT(log *logger.Logger) *JWT {
	return &JWT{log: log}
}

func (jwt *JWT) Create(u *user.User, duration time.Duration) (user.Token, error) {
	accessToken, err := token.New(data.User{
		Id:       u.ID,
		CIF:      u.CIF,
		FullName: u.FullName,
		Email:    u.Email,
	}, duration)
	if err != nil {
		jwt.log.Errorf("Create failed creating user (%v) token: %v", u.ID, err)
		return user.Token{}, err
	}

	jwt.log.Infof("Create success create user (%v) token", u.ID)
	return user.Token{
		AccessToken: accessToken.AccessToken,
		ExpiredTime: accessToken.ExpiredTime,
	}, nil
}
