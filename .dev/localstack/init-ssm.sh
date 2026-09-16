#!/usr/bin/env bash
# Seeds SSM Parameter Store with the config values expected by
# internal/common/config/config.go, under /go-api-template/local (matches
# ApplicationName + default APP_PROFILE="local").  Runs automatically on
# LocalStack startup (mounted into /etc/localstack/init/ready.d/).
set -euo pipefail

# LocalStack scopes SSM parameters by region, so this must match AWS_REGION
# used by the app (see Taskfile.yml dev task), not just DEFAULT_REGION (which
# is deprecated/ignored by LocalStack).
export AWS_DEFAULT_REGION="eu-west-1"

PATH_PREFIX="/go-api-template/local"

put() {
  awslocal ssm put-parameter --name "$PATH_PREFIX/$1" --type SecureString --value "$2" --overwrite
}

put "server_port" "8080"
put "server_debug" "true"
put "mongodb_uri" "mongodb://localhost:27017"

echo "Seeded SSM parameters under $PATH_PREFIX"
