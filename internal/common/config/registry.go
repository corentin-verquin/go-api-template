package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Registry struct {
	values map[string]any
}

func NewRegistry(logger *zap.Logger, params []Config) *Registry {
	appProfile := os.Getenv("APP_PROFILE")
	if appProfile == "" {
		appProfile = "local"
	}

	ssmPath := fmt.Sprintf("/%s/%s", ApplicationName, appProfile)
	ssmValues, err := loadFromSSM(context.Background(), ssmPath)
	if err != nil {
		logger.Sugar().Panicf("failed to load config from AWS SSM: %v", err)
	}

	values := make(map[string]any)
	for _, p := range params {
		key := strings.ToUpper(p.Key)
		rawValue := lo.CoalesceOrEmpty(os.Getenv(key), ssmValues[key])

		parsedValue := lo.TernaryF(
			rawValue == "",
			func() any { return p.DefaultValue },
			func() any {
				value, err := parse(rawValue, p.Type)
				if err != nil {
					logger.Sugar().Panicf("failed to parse config value for key: %s, error: %v", p.Key, err)
				}
				return value
			},
		)
		values[p.Key] = parsedValue
	}
	return &Registry{
		values: values,
	}
}

func parse(raw string, t configType) (any, error) {
	switch t {
	case Int:
		return strconv.Atoi(raw)
	case Bool:
		return strconv.ParseBool(raw)
	}
	return raw, nil
}

func (r *Registry) String(key string) string { return r.values[key].(string) }
func (r *Registry) Int(key string) int       { return r.values[key].(int) }
func (r *Registry) Bool(key string) bool     { return r.values[key].(bool) }
