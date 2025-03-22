//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/pkg/framework"
	"go.bankyaya.org/app/backend/pkg/service"
	"go.bankyaya.org/app/backend/pkg/util"
	"go.bankyaya.org/app/backend/pkg/util/config"
)

func initApp(cfg *config.Config) *app {
	wire.Build(
		framework.ProviderSet,
		service.ProviderSet,
		util.ProviderSet,
		echo.New,
		newApp,
	)
	return &app{}
}
