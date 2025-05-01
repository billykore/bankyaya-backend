//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter"
	"go.bankyaya.org/app/backend/internal/domain"
	"go.bankyaya.org/app/backend/internal/pkg"
	"go.bankyaya.org/app/backend/internal/pkg/config"
)

func initApp(cfg *config.Configs) *app {
	wire.Build(
		adapter.ProviderSet,
		domain.ProviderSet,
		pkg.ProviderSet,
		echo.New,
		newApp,
	)
	return &app{}
}
