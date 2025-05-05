package domain

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"go.bankyaya.org/app/backend/internal/domain/otp"
	"go.bankyaya.org/app/backend/internal/domain/user"
)

var ProviderSet = wire.NewSet(
	intrabank.NewService,
	user.NewService,
	otp.NewService,
)
