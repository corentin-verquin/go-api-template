package config

import (
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// loadFromSSM fetches every parameter under path (recursively) from AWS SSM
// Parameter Store, keyed by the upper-cased last segment of its name
func loadFromSSM(ctx context.Context, path string) (map[string]string, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := ssm.NewFromConfig(cfg)

	trueVal := true
	values := map[string]string{}
	var nextToken *string
	for {
		out, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           &path,
			Recursive:      &trueVal,
			WithDecryption: &trueVal,
			NextToken:      nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch parameters from SSM path %q: %w", path, err)
		}
		for _, p := range out.Parameters {
			segments := strings.Split(strings.Trim(*p.Name, "/"), "/")
			key := strings.ToUpper(segments[len(segments)-1])
			values[key] = *p.Value
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return values, nil
}
