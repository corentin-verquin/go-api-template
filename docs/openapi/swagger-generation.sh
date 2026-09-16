#!/bin/sh -e

go install github.com/swaggo/swag/v2/cmd/swag@latest

rm -f ./docs.go
rm -f ./swagger.yaml
rm -f ./swagger.json

swag init \
  --dir ../../ \
  --output ./ \
  --parseDependency \
  --generalInfo docs/openapi/api_info.go \
  --v3.1
