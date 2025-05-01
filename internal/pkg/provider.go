package pkg

import (
	"github.com/google/wire"
	"go.bankyaya.org/app/backend/internal/pkg/corebanking"
	"go.bankyaya.org/app/backend/internal/pkg/db/postgres"
	"go.bankyaya.org/app/backend/internal/pkg/email/mailtrap"
	"go.bankyaya.org/app/backend/internal/pkg/httpclient"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/notification/firebase"
	"go.bankyaya.org/app/backend/internal/pkg/validation"
)

var ProviderSet = wire.NewSet(
	logger.New,
	validation.New,
	postgres.New,
	corebanking.NewClient,
	httpclient.New,
	mailtrap.NewClient,
	firebase.New,
)
