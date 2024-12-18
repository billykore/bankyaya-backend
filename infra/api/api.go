package api

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/infra/api/corebanking"
)

var ProviderSet = wire.NewSet(
	corebanking.NewTransfer, wire.Bind(new(transfer.CoreBankingService), new(*corebanking.Transfer)),
)
