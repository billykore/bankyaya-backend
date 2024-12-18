//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain"
	"go.bankyaya.org/app/backend/infra/api"
	"go.bankyaya.org/app/backend/infra/email"
	"go.bankyaya.org/app/backend/infra/http"
	"go.bankyaya.org/app/backend/infra/storage"
	"go.bankyaya.org/app/backend/pkg"
	"go.bankyaya.org/app/backend/pkg/config"
)

func initApp(cfg *config.Config) *app {
	wire.Build(
		domain.ProviderSet,
		storage.ProviderSet,
		api.ProviderSet,
		email.ProviderSet,
		http.ProviderSet,
		pkg.ProviderSet,
		echo.New,
		newApp,
	)
	return &app{}
}
