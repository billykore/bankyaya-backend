package email

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/infra/email/mailer"
)

var ProviderSet = wire.NewSet(
	mailer.NewTransferEmail, wire.Bind(new(transfer.Email), new(*mailer.TransferEmail)),
)
