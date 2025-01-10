package http

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/infra/http/handler"
	"go.bankyaya.org/app/backend/infra/http/server"
)

var ProviderSet = wire.NewSet(
	handler.NewTransferHandler,
	handler.NewQRISHandler,
	handler.NewUserHandler,
	server.NewRouter,
	server.New,
)
