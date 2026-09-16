package admin

import (
	"net/http"
	"strings"

	"github.com/corentin-verquin/go-api-template/internal/common/auth"
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/gin-gonic/gin"
)

func AdminAuth(cfg *config.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := auth.Username(c.Request.Context())

		for _, cred := range strings.Split(cfg.String(adminCredentialsKey), ",") {
			if username == cred {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
