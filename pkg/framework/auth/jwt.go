package auth

import (
	"time"

	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/util/data"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	"go.bankyaya.org/app/backend/pkg/util/security/token"
)

type JWT struct {
	log *logger.Logger
}

func NewJWT(log *logger.Logger) *JWT {
	return &JWT{log: log}
}

func (jwt *JWT) Create(u *entity.User, duration time.Duration) (entity.Token, error) {
	accessToken, err := token.New(data.User{
		Id:       u.ID,
		CIF:      u.CIF,
		FullName: u.FullName,
		Email:    u.Email,
	}, duration)
	if err != nil {
		jwt.log.Errorf("Create failed creating user (%v) token: %v", u.ID, err)
		return entity.Token{}, err
	}

	jwt.log.Infof("Create success create user (%v) token", u.ID)
	return entity.Token{
		AccessToken: accessToken.AccessToken,
		ExpiredTime: accessToken.ExpiredTime,
	}, nil
}
