package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
)

// jwtConfig contains configuration for JWT auth middleware.
var jwtConfig = echojwt.Config{
	ContextKey:     ctxt.UserContextKey,
	SigningKey:     []byte(config.Load().Token.Secret),
	SuccessHandler: successHandler,
	ErrorHandler:   errorHandler,
}

// AuthenticateUser returns middleware function that validates token from headers
// and extract user information.
func AuthenticateUser() echo.MiddlewareFunc {
	return echojwt.WithConfig(jwtConfig)
}

// successHandler extract user information from token
// and save the information in the request context.
func successHandler(ctx echo.Context) {
	t := ctx.Get(ctxt.UserContextKey).(*jwt.Token)
	user := userFromToken(t)
	c := ctx.Request().Context()
	c = ctxt.ContextWithUser(c, &user)
	ctx.SetRequest(ctx.Request().WithContext(c))
}

// userFromToken returns user information from JWT token.
func userFromToken(token *jwt.Token) ctxt.User {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctxt.User{}
	}
	cif, ok := claims["cif"].(string)
	if !ok {
		return ctxt.User{}
	}
	userId, ok := claims["userId"].(int)
	if !ok {
		return ctxt.User{}
	}
	fullName, ok := claims["sub"].(string)
	if !ok {
		return ctxt.User{}
	}
	email, ok := claims["email"].(string)
	if !ok {
		return ctxt.User{}
	}
	return ctxt.User{
		CIF:      cif,
		Id:       userId,
		FullName: fullName,
		Email:    email,
	}
}

// errorHandler returns unauthorized response if there is an authentication error.
func errorHandler(ctx echo.Context, err error) error {
	return ctx.JSON(response.Unauthorized(err))
}
