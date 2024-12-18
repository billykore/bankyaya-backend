package domain

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/domain/transfer"
)

var ProviderSet = wire.NewSet(
	transfer.NewService,
)
