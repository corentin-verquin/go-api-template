package config

import (
	"os"

	"github.com/samber/lo"
)

var (
	ApplicationName      = "go-api-template"
	ApplicationVersion   = "dev"
	ApplicationGitHash   = ""
	ApplicationBuildDate = ""

	DeploymentDate, _ = lo.Coalesce(os.Getenv("DEPLOYMENT_DATE"), "unknown")
	Env, _            = lo.Coalesce(os.Getenv("ENV"), "local")
)

const (
	BaseURI = "/api/go-api-template"
)
