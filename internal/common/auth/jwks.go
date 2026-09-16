package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewJWKS(lc fx.Lifecycle, cfg *config.Registry, log *zap.Logger) (*keyfunc.JWKS, error) {
	options := keyfunc.Options{
		RefreshInterval:   time.Hour * 12,
		RefreshTimeout:    time.Second * 5,
		RefreshUnknownKID: true,
		RefreshErrorHandler: func(err error) {
			log.Error("failed to refresh JWKS", zap.Error(err))
		},
	}
	cognitoJwksURL := fmt.Sprintf("%s/%s/.well-known/jwks.json", cfg.String(cognitoUrl), cfg.String(cognitoUserPoolId))

	jwks, err := keyfunc.Get(cognitoJwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWKS from %s: %w", cognitoJwksURL, err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			jwks.EndBackground()
			return nil
		},
	})

	return jwks, nil
}
