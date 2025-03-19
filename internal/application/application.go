package application

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/application/qris"
	"go.bankyaya.org/app/backend/internal/application/scheduler"
	"go.bankyaya.org/app/backend/internal/application/transfer"
	"go.bankyaya.org/app/backend/internal/application/user"
)

var ProviderSet = wire.NewSet(
	transfer.NewUsecase,
	user.NewUsecase,
	qris.NewUsecase,
	scheduler.NewUsecase,
)
