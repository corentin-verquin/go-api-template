package sample

import (
	"github.com/corentin-verquin/go-api-template/internal/common/logger"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) SampleMethod() string {
	done := logger.FeatureLogger(s.logger, Pkg, "SampleMethod")
	defer done()
	return "Hello from service"
}
