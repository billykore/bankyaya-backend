package framework

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/pkg/framework/auth"
	"go.bankyaya.org/app/backend/pkg/framework/corebanking"
	"go.bankyaya.org/app/backend/pkg/framework/email"
	"go.bankyaya.org/app/backend/pkg/framework/http/handler"
	"go.bankyaya.org/app/backend/pkg/framework/http/server"
	"go.bankyaya.org/app/backend/pkg/framework/messaging"
	"go.bankyaya.org/app/backend/pkg/framework/qris"
	"go.bankyaya.org/app/backend/pkg/framework/storage/repo"
	"go.bankyaya.org/app/backend/pkg/interface/api"
	emailinterface "go.bankyaya.org/app/backend/pkg/interface/email"
	messaginginterface "go.bankyaya.org/app/backend/pkg/interface/messaging"
	"go.bankyaya.org/app/backend/pkg/interface/repository"
	"go.bankyaya.org/app/backend/pkg/interface/security"
)

var authProviderSet = wire.NewSet(
	auth.NewJWT, wire.Bind(new(security.TokenService), new(*auth.JWT)),
	auth.NewBcryptHasher, wire.Bind(new(security.PasswordHasher), new(*auth.BcryptHasher)),
)

var coreBankingProviderSet = wire.NewSet(
	corebanking.New, wire.Bind(new(api.CoreBanking), new(*corebanking.CoreBanking)),
)

var emailProviderSet = wire.NewSet(
	email.NewQRISEmail, wire.Bind(new(emailinterface.QRISReceiptMailer), new(*email.QRISEmail)),
	email.NewTransferEmail, wire.Bind(new(emailinterface.TransferReceiptMailer), new(*email.TransferEmail)),
)

var messagingProviderSet = wire.NewSet(
	messaging.NewSchedulerPublisher, wire.Bind(new(messaginginterface.AutoDebitEventPublisher), new(*messaging.SchedulerPublisher)),
	messaging.NewTransferConsumer,
	messaging.NewListener,
)

var qrisProviderSet = wire.NewSet(
	qris.NewQRIS, wire.Bind(new(api.QRIS), new(*qris.QRIS)),
)

var repositoryProviderSet = wire.NewSet(
	repo.NewTransferRepo, wire.Bind(new(repository.TransferRepository), new(*repo.TransferRepo)),
	repo.NewUserRepo, wire.Bind(new(repository.UserRepository), new(*repo.UserRepo)),
	repo.NewSchedulerRepo, wire.Bind(new(repository.ScheduleRepository), new(*repo.SchedulerRepo)),
)

var handlerProviderSet = wire.NewSet(
	handler.NewTransfer,
	handler.NewQRIS,
	handler.NewUserHandler,
	handler.NewScheduler,
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
