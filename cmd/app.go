package cmd

import (
	"github.com/corentin-verquin/go-api-template/internal/common/admin"
	"github.com/corentin-verquin/go-api-template/internal/common/auth"
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/corentin-verquin/go-api-template/internal/common/ginhttp"
	"github.com/corentin-verquin/go-api-template/internal/common/logger"
	mongodb "github.com/corentin-verquin/go-api-template/internal/common/mongoDB"
	"github.com/corentin-verquin/go-api-template/internal/sample"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Execute() {
	deps := fx.Options(
		logger.Module,
		config.Module,
		auth.Module,
		admin.Module,
		mongodb.Module,
	)

	feature := fx.Options(
		sample.Module,
	)

	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			if log.Level() == zap.DebugLevel {
				return &fxevent.ZapLogger{Logger: log}
			}
			return &fxevent.NopLogger
		}),

		deps,
		feature,

		ginhttp.Module,
		ginhttp.Start(),
	).Run()
}
