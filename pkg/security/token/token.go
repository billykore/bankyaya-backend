package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.bankyaya.org/app/backend/pkg/config"
	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/uuid"
)

const tokenExpiredTime = 15 * time.Minute

// Token contains access token and expired time of the token
type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiredTime int64  `json:"expiredTime"`
}

// New return new generated token for the given username.
func New(subject string) (Token, error) {
	cfg := config.Get()
	return generateToken(cfg, subject)
}

func generateToken(cfg *config.Config, subject string) (Token, error) {
	id, err := uuid.New()
	if err != nil {
		return Token{}, err
	}
	now := time.Now()
	exp := now.Add(tokenExpiredTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": id,
		"sub": subject,
		"iss": "app-name",
		"aud": subject,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	token.Header["kid"] = cfg.Token.HeaderKid

	t, err := token.SignedString([]byte(cfg.Token.Secret))
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken: t,
		ExpiredTime: exp.Unix(),
	}, nil
}

// UserFromToken returns user information from JWT token.
func UserFromToken(token *jwt.Token) entity.User {
	claims := token.Claims.(jwt.MapClaims)
	return entity.User{
		Id: claims["userId"].(int),
	}
}
