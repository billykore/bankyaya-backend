package http

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/infra/http/handler"
	"go.bankyaya.org/app/backend/infra/http/server"
)

var ProviderSet = wire.NewSet(
	handler.NewTransferHandler,
	server.NewRouter,
	server.New,
)