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
	"go.bankyaya.org/app/backend/internal/adapter/sequence"
	"go.bankyaya.org/app/backend/internal/adapter/storage/repo"
	"go.bankyaya.org/app/backend/internal/core/port/api"
	emailport "go.bankyaya.org/app/backend/internal/core/port/email"
	messagingport "go.bankyaya.org/app/backend/internal/core/port/messaging"
	repositoryport "go.bankyaya.org/app/backend/internal/core/port/repository"
	"go.bankyaya.org/app/backend/internal/core/port/security"
)

var authProviderSet = wire.NewSet(
	auth.NewJWT, wire.Bind(new(security.TokenService), new(*auth.JWT)),
	auth.NewBcryptHasher, wire.Bind(new(security.PasswordHasher), new(*auth.BcryptHasher)),
)

var coreBankingProviderSet = wire.NewSet(
	corebanking.New, wire.Bind(new(api.CoreBanking), new(*corebanking.CoreBanking)),
)

var emailProviderSet = wire.NewSet(
	email.NewQRISEmail, wire.Bind(new(emailport.QRISReceiptMailer), new(*email.QRISEmail)),
	email.NewTransferEmail, wire.Bind(new(emailport.TransferReceiptMailer), new(*email.TransferEmail)),
)

var messagingProviderSet = wire.NewSet(
	messaging.NewSchedulerPublisher, wire.Bind(new(messagingport.AutoDebitEventPublisher), new(*messaging.SchedulerPublisher)),
	messaging.NewTransferConsumer,
	messaging.NewListener,
)

var qrisProviderSet = wire.NewSet(
	qris.NewQRIS, wire.Bind(new(api.QRIS), new(*qris.QRIS)),
)

var sequencerProviderSet = wire.NewSet(
	sequence.New, wire.Bind(new(security.SequenceGenerator), new(*sequence.Sequence)),
)

var repositoryProviderSet = wire.NewSet(
	repo.NewTransferRepo, wire.Bind(new(repositoryport.TransferRepository), new(*repo.TransferRepo)),
	repo.NewUserRepo, wire.Bind(new(repositoryport.UserRepository), new(*repo.UserRepo)),
	repo.NewSchedulerRepo, wire.Bind(new(repositoryport.ScheduleRepository), new(*repo.SchedulerRepo)),
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
	sequencerProviderSet,
	repositoryProviderSet,
	handlerProviderSet,
	serverProviderSet,
)
