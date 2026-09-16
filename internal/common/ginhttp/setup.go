package ginhttp

import (
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const (
	Pkg = "gin_http"

	serverPort  = "server_port"
	serverDebug = "server_debug"
)

func SetupParams() (params []config.Config) {
	return []config.Config{
		config.NewInt(serverPort, 8080),
		config.NewBool(serverDebug, false),
	}
}

func Start() fx.Option {
	return fx.Invoke(func(*gin.Engine) {})
}

var Module = fx.Module(Pkg,
	config.AsConfigParams(SetupParams),
	fx.Provide(NewServer),
	fx.Provide(NewRouter),
)
