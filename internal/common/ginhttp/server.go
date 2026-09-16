package ginhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/corentin-verquin/go-api-template/internal/common/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer(lc fx.Lifecycle, cfg *config.Registry, log *zap.Logger) *gin.Engine {
	if cfg.Bool(serverDebug) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(logger.ZapGinLogger(log), logger.ZapGinRecovery(log))
	router.GET("/_info", info)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Int(serverPort)),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Error("fail to start HTTP Server", zap.Error(err))
				return err
			}

			for _, route := range router.Routes() {
				log.Info("HTTP route registered", zap.String("method", route.Method), zap.String("path", route.Path))
			}

			go srv.Serve(ln)
			log.Info("succeeded to start HTTP Server", zap.String("addr", srv.Addr))
			return nil

		},
		OnStop: func(ctx context.Context) error {
			log.Info("HTTP Server is stopping")
			return srv.Shutdown(ctx)
		},
	})

	return router
}

func info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":            config.ApplicationName,
		"version":         config.ApplicationVersion,
		"commit_hash":     config.ApplicationGitHash,
		"build_date":      config.ApplicationBuildDate,
		"env":             config.Env,
		"deployment_date": config.DeploymentDate,
	})
}
