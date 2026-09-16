#!/usr/bin/env bash
# Authenticates the test user against cognito-local and prints an ID token,
# ready to use as `Authorization: Bearer <token>` against the API.
# Requires ./init-cognito.sh (or `task cognito:init`) to have been run first.
# Usage: ./dev/cognito-local/generate-jwt.sh (or `task generate:jwt`)
set -euo pipefail

ENDPOINT="http://localhost:9229"
export AWS_ACCESS_KEY_ID="local"
export AWS_SECRET_ACCESS_KEY="local"
export AWS_DEFAULT_REGION="local"

POOL_NAME="template-api-local"
CLIENT_NAME="template-api-client"
TEST_USERNAME="test@template.local"
TEST_PASSWORD="Test1234!"

cognito() {
  aws --no-cli-pager --endpoint-url "$ENDPOINT" cognito-idp "$@"
}

USER_POOL_ID=$(cognito list-user-pools --max-results 60 \
  --query "UserPools[?Name=='$POOL_NAME'].Id | [0]" --output text)

if [[ -z "$USER_POOL_ID" || "$USER_POOL_ID" == "None" ]]; then
  echo "User pool '$POOL_NAME' not found, run 'task cognito:init' first." >&2
  exit 1
fi

CLIENT_ID=$(cognito list-user-pool-clients --user-pool-id "$USER_POOL_ID" \
  --query "UserPoolClients[?ClientName=='$CLIENT_NAME'].ClientId | [0]" --output text)

if [[ -z "$CLIENT_ID" || "$CLIENT_ID" == "None" ]]; then
  echo "App client '$CLIENT_NAME' not found, run 'task cognito:init' first." >&2
  exit 1
fi

ID_TOKEN=$(cognito initiate-auth \
  --client-id "$CLIENT_ID" \
  --auth-flow USER_PASSWORD_AUTH \
  --auth-parameters USERNAME="$TEST_USERNAME",PASSWORD="$TEST_PASSWORD" \
  --query "AuthenticationResult.AccessToken" --output text)

echo "$ID_TOKEN"
