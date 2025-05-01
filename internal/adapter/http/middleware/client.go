package middleware

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/version"
)

const (
	clientNameHeaderKey    = "X-Client-Name"
	clientVersionHeaderKey = "X-Client-Version"
)

var errInvalidClient = errors.New("invalid client")

// ValidateClients returns an Echo middleware that validates client name and version
// from the request headers against a list of allowed clients defined in the configuration.
func ValidateClients() echo.MiddlewareFunc {
	cfg := config.Load()
	clients := cfg.Clients

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			clientFromHeader := ctx.Request().Header.Get(clientNameHeaderKey)
			versionFromHeader := ctx.Request().Header.Get(clientVersionHeaderKey)

			for _, client := range clients {
				if clientFromHeader == client.Name {
					minVersion := version.NewVersion(client.MinVersion)
					maxVersion := version.NewVersion(client.MaxVersion)
					headerSemanticVersion := version.NewVersion(versionFromHeader)

					if !headerSemanticVersion.Between(minVersion, maxVersion) {
						return ctx.JSON(response.Forbidden(errInvalidClient))
					}
					return next(ctx)
				}
			}

			return ctx.JSON(response.Forbidden(errInvalidClient))
		}
	}
}
