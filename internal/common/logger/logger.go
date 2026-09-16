package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	zapConfig := lo.Ternary(config.Env == "local", zap.NewDevelopmentConfig(), zap.NewProductionConfig())
	zapConfig.InitialFields = map[string]any{
		"env":     config.Env,
		"service": config.ApplicationName,
	}
	if config.Env == "local" {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func ZapGinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func ZapGinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", zap.Any("error", err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func FeatureLogger(logger *zap.Logger, pkg string, featureName string, fields ...zap.Field) func() {
	start := time.Now()
	logger.Info(
		fmt.Sprintf("[%s] %s started", pkg, featureName),
		fields...,
	)
	return func() {
		logger.Info(
			fmt.Sprintf("[%s] %s completed", pkg, featureName),
			append(fields, zap.Duration("duration", time.Since(start)))...,
		)
	}
}
