package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.bankyaya.org/app/backend/pkg/util/config"
	"go.bankyaya.org/app/backend/pkg/util/data"
	"go.bankyaya.org/app/backend/pkg/util/uuid"
)

// Token contains access token and expired time of the token
type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiredTime int64  `json:"expiredTime"`
}

// New return new generated token for the given username.
func New(user data.User, duration time.Duration) (Token, error) {
	cfg := config.Get()
	return generateToken(cfg, user, duration)
}

func generateToken(cfg *config.Config, user data.User, duration time.Duration) (Token, error) {
	id, err := uuid.New()
	if err != nil {
		return Token{}, err
	}
	now := time.Now()
	exp := now.Add(duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":    id,
		"sub":    user.FullName,
		"iss":    "api.bankyaya.co.id",
		"aud":    "https://bankyaya.co.id",
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
		"cif":    user.CIF,
		"userId": user.Id,
		"email":  user.Email,
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
func UserFromToken(token *jwt.Token) data.User {
	claims := token.Claims.(jwt.MapClaims)
	cif, ok := claims["cif"].(string)
	if !ok {
		return data.User{}
	}
	userId, ok := claims["userId"].(int)
	if !ok {
		return data.User{}
	}
	fullName, ok := claims["sub"].(string)
	if !ok {
		return data.User{}
	}
	email, ok := claims["email"].(string)
	if !ok {
		return data.User{}
	}
	return data.User{
		CIF:      cif,
		Id:       userId,
		FullName: fullName,
		Email:    email,
	}
}
