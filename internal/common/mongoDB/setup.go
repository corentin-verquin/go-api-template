package mongodb

import (
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"go.uber.org/fx"
)

const (
	Pkg = "mongodb"

	mongoDBUriKey = "mongodb_uri"
)

func SetupParams() (params []config.Config) {
	return []config.Config{
		config.NewString(mongoDBUriKey, ""),
	}
}

var Module = fx.Module(Pkg,
	config.AsConfigParams(SetupParams),
	fx.Provide(NewMongoDB),
	fx.Invoke(func(lc fx.Lifecycle, db *MongoDB) {
		lc.Append(fx.Hook{
			OnStart: db.onStart,
			OnStop:  db.onStop,
		})
	}),
)
