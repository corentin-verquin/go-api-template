package openapi

//go:generate ./swagger-generation.sh

// @title			API go template

// @securityDefinitions.apikey	Cognito
// @in							header
// @name						Authorization
// @description				    Cognito JWT (use "Bearer XYZ" as value)
