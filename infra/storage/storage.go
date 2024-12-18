package storage

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/infra/storage/repo"
)

var ProviderSet = wire.NewSet(
	repo.NewTransferRepo, wire.Bind(new(transfer.Repository), new(*repo.TransferRepo)),
)
