package pkg

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/pkg/corebanking"
	"go.bankyaya.org/app/backend/pkg/db/postgres"
	"go.bankyaya.org/app/backend/pkg/email/mailtrap"
	"go.bankyaya.org/app/backend/pkg/httpclient"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/qris"
	"go.bankyaya.org/app/backend/pkg/validation"
)

var ProviderSet = wire.NewSet(
	logger.New,
	validation.New,
	postgres.New,
	corebanking.NewClient,
	httpclient.New,
	mailtrap.NewClient,
	qris.NewClient,
)
