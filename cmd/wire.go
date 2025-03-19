//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter"
	"go.bankyaya.org/app/backend/internal/application"
	"go.bankyaya.org/app/backend/internal/core"
	"go.bankyaya.org/app/backend/pkg"
	"go.bankyaya.org/app/backend/pkg/config"
)

func initApp(cfg *config.Config) *app {
	wire.Build(
		adapter.ProviderSet,
		application.ProviderSet,
		core.ProviderSet,
		pkg.ProviderSet,
		echo.New,
		newApp,
	)
	return &app{}
}
