package auth

import (
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"go.uber.org/fx"
)

const (
	Pkg = "auth"

	cognitoUrl        = "cognito_url"
	cognitoClientId   = "cognito_client_id"
	cognitoUserPoolId = "cognito_user_pool_id"
)

func SetupParams() (params []config.Config) {
	return []config.Config{
		config.NewString(cognitoUrl, ""),
		config.NewString(cognitoClientId, ""),
		config.NewString(cognitoUserPoolId, ""),
	}
}

var Module = fx.Module(Pkg,
	config.AsConfigParams(SetupParams),
	fx.Provide(NewJWKS),
)
