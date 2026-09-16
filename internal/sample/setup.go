package sample

import (
	"github.com/corentin-verquin/go-api-template/internal/common/ginhttp"
	"go.uber.org/fx"
)

const Pkg = "sample"

func SetupRoutes(h *Handler, r *ginhttp.Router) {
	r.Public.GET("/hello", h.SampleHandler)
	r.Secured.GET("/secure", h.SampleHandler)
	r.Admin.GET("/admin", h.SampleHandler)
}

var Module = fx.Module(Pkg,
	fx.Provide(NewService),
	fx.Provide(NewHandler),
	fx.Invoke(SetupRoutes),
)
