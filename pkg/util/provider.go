package util

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/pkg/util/corebanking"
	"go.bankyaya.org/app/backend/pkg/util/db/postgres"
	"go.bankyaya.org/app/backend/pkg/util/email/mailtrap"
	"go.bankyaya.org/app/backend/pkg/util/httpclient"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	"go.bankyaya.org/app/backend/pkg/util/messaging/rabbitmq"
	"go.bankyaya.org/app/backend/pkg/util/qris"
	"go.bankyaya.org/app/backend/pkg/util/validation"
)

var ProviderSet = wire.NewSet(
	logger.New,
	validation.New,
	postgres.New,
	corebanking.NewClient,
	httpclient.New,
	mailtrap.NewClient,
	qris.NewClient,
	rabbitmq.NewConnection,
)
