package core

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/core/service"
)

var ProviderSet = wire.NewSet(
	service.NewQRIS,
	service.NewScheduler,
	service.NewTransfer,
	service.NewUser,
)
