package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.bankyaya.org/app/backend/internal/domain/user"
	"go.bankyaya.org/app/backend/internal/pkg/config"
)

type JWT struct {
	cfg *config.Configs
}

func NewJWT(cfg *config.Configs) *JWT {
	return &JWT{
		cfg: cfg,
	}
}

func (j *JWT) Create(u *user.User, duration time.Duration) (*user.Token, error) {
	now := time.Now()
	exp := now.Add(duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":    uuid.New(),
		"sub":    u.FullName,
		"iss":    "api.bankyaya.co.id",
		"aud":    "https://bankyaya.co.id",
		"exp":    exp.Unix(),
		"iat":    now.Unix(),
		"cif":    u.CIF,
		"userId": u.ID,
		"email":  u.Email,
	})
	token.Header["kid"] = j.cfg.Token.HeaderKid

	tokenString, err := token.SignedString([]byte(j.cfg.Token.Secret))
	if err != nil {
		return nil, err
	}

	return &user.Token{
		AccessToken: tokenString,
		ExpiresAt:   exp,
	}, nil
}
