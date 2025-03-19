package core

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/core/qris"
	"go.bankyaya.org/app/backend/internal/core/scheduler"
	"go.bankyaya.org/app/backend/internal/core/transfer"
	"go.bankyaya.org/app/backend/internal/core/user"
)

var ProviderSet = wire.NewSet(
	transfer.NewService,
	qris.NewService,
	user.NewService,
	scheduler.NewService,
)
