package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type cognitoClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	ClientID string `json:"client_id"`
	TokenUse string `json:"token_use"`

	expectedClientID string
}

func (c cognitoClaims) Validate() error {
	if c.TokenUse != "access" {
		return fmt.Errorf("unexpected token_use %q", c.TokenUse)
	}
	if c.ClientID != c.expectedClientID {
		return errors.New("client_id mismatch")
	}
	if c.Username == "" {
		return errors.New("missing username claim")
	}
	return nil
}

func JwtAuth(jwks *keyfunc.JWKS, cfg *config.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		cognitoURL := fmt.Sprintf("%s/%s", cfg.String(cognitoUrl), cfg.String(cognitoUserPoolId))

		claims := cognitoClaims{expectedClientID: cfg.String(cognitoClientId)}
		token, err := jwt.ParseWithClaims(
			authHeader[7:],
			&claims,
			jwks.Keyfunc,
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(cognitoURL),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Request = c.Request.WithContext(WithUsername(ctx, claims.Username))
		c.Next()
	}
}
