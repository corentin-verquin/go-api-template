package admin

import (
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"go.uber.org/fx"
)

const (
	Pkg = "admin"

	adminCredentialsKey = "admin_credentials"
)

func SetupParams() (params []config.Config) {
	params = []config.Config{
		config.NewString(adminCredentialsKey, ""),
	}
	return
}

var Module = fx.Module(Pkg, config.AsConfigParams(SetupParams))
