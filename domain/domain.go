package domain

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/domain/qris"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/domain/user"
)

var ProviderSet = wire.NewSet(
	transfer.NewService,
	qris.NewService,
	user.NewService,
)
