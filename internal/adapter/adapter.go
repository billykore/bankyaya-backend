package adapter

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/adapter/auth"
	"go.bankyaya.org/app/backend/internal/adapter/corebanking"
	"go.bankyaya.org/app/backend/internal/adapter/email"
	"go.bankyaya.org/app/backend/internal/adapter/http/handler"
	"go.bankyaya.org/app/backend/internal/adapter/http/server"
	"go.bankyaya.org/app/backend/internal/adapter/messaging"
	"go.bankyaya.org/app/backend/internal/adapter/qris"
	"go.bankyaya.org/app/backend/internal/adapter/storage/repo"
	qriscore "go.bankyaya.org/app/backend/internal/core/qris"
	"go.bankyaya.org/app/backend/internal/core/scheduler"
	"go.bankyaya.org/app/backend/internal/core/transfer"
	"go.bankyaya.org/app/backend/internal/core/user"
)

var authProviderSet = wire.NewSet(
	auth.NewJWT, wire.Bind(new(user.TokenService), new(*auth.JWT)),
	auth.NewBcryptHasher, wire.Bind(new(user.PasswordHasher), new(*auth.BcryptHasher)),
)

var coreBankingProviderSet = wire.NewSet(
	corebanking.NewTransfer, wire.Bind(new(transfer.CoreBanking), new(*corebanking.Transfer)),
	corebanking.NewQRIS, wire.Bind(new(qriscore.CoreBanking), new(*corebanking.QRIS)),
)

var emailProviderSet = wire.NewSet(
	email.NewQRISEmail, wire.Bind(new(qriscore.ReceiptMailer), new(*email.QRISEmail)),
	email.NewTransferEmail, wire.Bind(new(transfer.ReceiptMailer), new(*email.TransferEmail)),
)

var messagingProviderSet = wire.NewSet(
	messaging.NewSchedulerPublisher, wire.Bind(new(scheduler.AutoDebitEventPublisher), new(*messaging.SchedulerPublisher)),
	messaging.NewTransferConsumer,
	messaging.NewListener,
)

var qrisProviderSet = wire.NewSet(
	qris.NewQRIS, wire.Bind(new(qriscore.QRIS), new(*qris.QRIS)),
)

var repositoryProviderSet = wire.NewSet(
	repo.NewTransferRepo, wire.Bind(new(transfer.Repository), new(*repo.TransferRepo)),
	repo.NewUserRepo, wire.Bind(new(user.Repository), new(*repo.UserRepo)),
	repo.NewSchedulerRepo, wire.Bind(new(scheduler.Repository), new(*repo.SchedulerRepo)),
)

var handlerProviderSet = wire.NewSet(
	handler.NewTransferHandler,
	handler.NewQRISHandler,
	handler.NewUserHandler,
	handler.NewSchedulerHandler,
)

var serverProviderSet = wire.NewSet(
	server.NewRouter,
	server.New,
)

var ProviderSet = wire.NewSet(
	authProviderSet,
	coreBankingProviderSet,
	emailProviderSet,
	messagingProviderSet,
	qrisProviderSet,
	repositoryProviderSet,
	handlerProviderSet,
	serverProviderSet,
)
