package config

import "go.uber.org/fx"

const Pkg = "config"

func AsConfigParams(fn func() []Config) fx.Option {
	return fx.Provide(
		fx.Annotate(
			fn,
			fx.ResultTags(`group:"config_params,flatten"`),
		),
	)
}

func ProvideRegistry() fx.Option {
	return fx.Provide(
		fx.Annotate(
			NewRegistry,
			fx.ParamTags(``, `group:"config_params"`),
		),
	)
}

var Module = fx.Module(Pkg, ProvideRegistry())
