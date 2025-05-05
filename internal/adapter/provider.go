package adapter

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/adapter/corebanking"
	"go.bankyaya.org/app/backend/internal/adapter/email"
	"go.bankyaya.org/app/backend/internal/adapter/http/handler"
	"go.bankyaya.org/app/backend/internal/adapter/http/server"
	"go.bankyaya.org/app/backend/internal/adapter/notification"
	"go.bankyaya.org/app/backend/internal/adapter/otp"
	"go.bankyaya.org/app/backend/internal/adapter/password"
	"go.bankyaya.org/app/backend/internal/adapter/sequence"
	"go.bankyaya.org/app/backend/internal/adapter/storage/repo"
	"go.bankyaya.org/app/backend/internal/adapter/token"
	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	otpdomain "go.bankyaya.org/app/backend/internal/domain/otp"
	"go.bankyaya.org/app/backend/internal/domain/user"
)

var tokenProviderSet = wire.NewSet(
	token.NewJWT, wire.Bind(new(user.TokenService), new(*token.JWT)),
)

var passwordProviderSet = wire.NewSet(
	password.NewBcryptHasher, wire.Bind(new(user.PasswordHasher), new(*password.BcryptHasher)))

var coreBankingProviderSet = wire.NewSet(
	corebanking.NewIntrabankCoreBanking, wire.Bind(new(intrabank.CoreBanking), new(*corebanking.IntrabankCoreBanking)),
)

var emailProviderSet = wire.NewSet(
	email.NewTransferEmail, wire.Bind(new(intrabank.ReceiptMailer), new(*email.IntrabankEmail)),
	email.NewOTPEmail, wire.Bind(new(otpdomain.Sender), new(*email.OTPEmail)),
)

var notificationProviderSet = wire.NewSet(
	notification.NewIntrabankNotification, wire.Bind(new(intrabank.Notifier), new(*notification.IntrabankNotification)),
)

var sequencerProviderSet = wire.NewSet(
	sequence.New, wire.Bind(new(intrabank.SequenceGenerator), new(*sequence.UUID)),
)

var otpProviderSet = wire.NewSet(
	otp.NewOTP, wire.Bind(new(otpdomain.Generator), new(*otp.OTP)),
)

var repositoryProviderSet = wire.NewSet(
	repo.NewIntrabankRepo, wire.Bind(new(intrabank.Repository), new(*repo.IntrabankRepo)),
	repo.NewUserRepo, wire.Bind(new(user.Repository), new(*repo.UserRepo)),
	repo.NewOTPRepo, wire.Bind(new(otpdomain.Repository), new(*repo.OTPRepo)),
)

var handlerProviderSet = wire.NewSet(
	handler.NewIntrabankHandler,
	handler.NewUserHandler,
	handler.NewOTPHandler,
)

var serverProviderSet = wire.NewSet(
	server.NewRouter,
	server.New,
)

var ProviderSet = wire.NewSet(
	tokenProviderSet,
	passwordProviderSet,
	coreBankingProviderSet,
	emailProviderSet,
	notificationProviderSet,
	sequencerProviderSet,
	otpProviderSet,
	repositoryProviderSet,
	handlerProviderSet,
	serverProviderSet,
)
