#!/usr/bin/env bash
# Provisions a User Pool + App Client + test user against cognito-local
# (docker-compose service "cognito-local", https://github.com/jagregory/cognito-local
# — used instead of LocalStack's cognito-idp, which is a Pro-only feature).
# Idempotent: safe to re-run, skips resources that already exist.
# Usage: ./dev/cognito-local/init-cognito.sh (or `task cognito:init`)
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
  aws --endpoint-url "$ENDPOINT" cognito-idp "$@"
}

for _ in $(seq 1 30); do
  if curl -sf -o /dev/null "$ENDPOINT"; then
    break
  fi
  sleep 1
done

USER_POOL_ID=$(cognito list-user-pools --max-results 60 \
  --query "UserPools[?Name=='$POOL_NAME'].Id | [0]" --output text)

if [[ -z "$USER_POOL_ID" || "$USER_POOL_ID" == "None" ]]; then
  USER_POOL_ID=$(cognito create-user-pool --pool-name "$POOL_NAME" --query "UserPool.Id" --output text)
  echo "Created user pool: $USER_POOL_ID"
else
  echo "User pool already exists: $USER_POOL_ID"
fi

CLIENT_ID=$(cognito list-user-pool-clients --user-pool-id "$USER_POOL_ID" \
  --query "UserPoolClients[?ClientName=='$CLIENT_NAME'].ClientId | [0]" --output text)

if [[ -z "$CLIENT_ID" || "$CLIENT_ID" == "None" ]]; then
  CLIENT_ID=$(cognito create-user-pool-client \
    --user-pool-id "$USER_POOL_ID" \
    --client-name "$CLIENT_NAME" \
    --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH \
    --query "UserPoolClient.ClientId" --output text)
  echo "Created app client: $CLIENT_ID"
else
  echo "App client already exists: $CLIENT_ID"
fi

if cognito admin-get-user --user-pool-id "$USER_POOL_ID" --username "$TEST_USERNAME" >/dev/null 2>&1; then
  echo "Test user already exists: $TEST_USERNAME"
else
  cognito admin-create-user \
    --user-pool-id "$USER_POOL_ID" \
    --username "$TEST_USERNAME" \
    --user-attributes Name=email,Value="$TEST_USERNAME" Name=email_verified,Value=true \
    --message-action SUPPRESS >/dev/null
  cognito admin-set-user-password \
    --user-pool-id "$USER_POOL_ID" \
    --username "$TEST_USERNAME" \
    --password "$TEST_PASSWORD" \
    --permanent
  echo "Created test user: $TEST_USERNAME / $TEST_PASSWORD"
fi

# Seed the generated IDs into the same LocalStack SSM path config.go reads (see init-ssm.sh).
SSM_ENDPOINT="http://localhost:4566"
SSM_PATH_PREFIX="/go-api-template/local"
ssm() {
  AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_DEFAULT_REGION="eu-west-1" \
    aws --endpoint-url "$SSM_ENDPOINT" ssm "$@"
}
ssm put-parameter --name "$SSM_PATH_PREFIX/cognito_user_pool_id" --type SecureString --value "$USER_POOL_ID" --overwrite >/dev/null
ssm put-parameter --name "$SSM_PATH_PREFIX/cognito_url" --type SecureString --value "$ENDPOINT" --overwrite >/dev/null
ssm put-parameter --name "$SSM_PATH_PREFIX/cognito_client_id" --type SecureString --value "$CLIENT_ID" --overwrite >/dev/null
echo "Seeded cognito_user_pool_id/cognito_client_id/cognito_url under $SSM_PATH_PREFIX"

echo "Cognito Local User Pool ID: $USER_POOL_ID"
echo "Cognito Local App Client ID: $CLIENT_ID"
