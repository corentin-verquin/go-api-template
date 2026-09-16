package ginhttp

import (
	"github.com/MicahParks/keyfunc/v2"
	"github.com/corentin-verquin/go-api-template/internal/common/admin"
	"github.com/corentin-verquin/go-api-template/internal/common/auth"
	"github.com/corentin-verquin/go-api-template/internal/common/config"
	"github.com/gin-gonic/gin"
)

type publicRouter struct{ *gin.RouterGroup }
type adminRouter struct{ *gin.RouterGroup }
type securedRouter struct{ *gin.RouterGroup }

func newPublicRouter(server *gin.Engine) publicRouter {
	return publicRouter{server.Group(config.BaseURI)}
}
func newSecuredRouter(public publicRouter, cfg *config.Registry, keys *keyfunc.JWKS) securedRouter {
	return securedRouter{public.Group("", auth.JwtAuth(keys, cfg))}
}
func newAdminRouter(secured securedRouter, cfg *config.Registry) adminRouter {
	return adminRouter{secured.Group("", admin.AdminAuth(cfg))}
}

type Router struct {
	Public  *publicRouter
	Secured *securedRouter
	Admin   *adminRouter
}

func NewRouter(server *gin.Engine, cfg *config.Registry, keys *keyfunc.JWKS) *Router {
	p := newPublicRouter(server)
	s := newSecuredRouter(p, cfg, keys)
	a := newAdminRouter(s, cfg)

	return &Router{
		Public:  &p,
		Secured: &s,
		Admin:   &a,
	}
}
