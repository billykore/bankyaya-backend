package api

import (
	"github.com/google/wire"
	qrisdomain "go.bankyaya.org/app/backend/domain/qris"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/infra/api/corebanking"
	"go.bankyaya.org/app/backend/infra/api/qris"
)

var ProviderSet = wire.NewSet(
	corebanking.NewTransfer, wire.Bind(new(transfer.CoreBanking), new(*corebanking.Transfer)),
	corebanking.NewQRIS, wire.Bind(new(qrisdomain.CoreBanking), new(*corebanking.QRIS)),
	qris.NewQRIS, wire.Bind(new(qrisdomain.QRIS), new(*qris.QRIS)),
)
