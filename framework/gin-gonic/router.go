package gin_gonic

import (
	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
)

// Router .
type Router struct {
	*gin.RouterGroup
}

// Handler 业务处理
type Handler func(ctx *hfwctx.Context)

// h .
func (h Handler) h() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h(hfwctx.Encode(ctx))
	}
}

// NewRouterGroup .
func (r *Router) NewRouterGroup(p string) *Router {
	return &Router{r.RouterGroup.Group(p)}
}

// GET .
func (r *Router) GET(uri string, h Handler) {
	r.RouterGroup.GET(uri, h.h())
}

// PUT .
func (r *Router) PUT(uri string, h Handler) {
	r.RouterGroup.PUT(uri, h.h())
}

// POST .
func (r *Router) POST(uri string, h Handler) {
	r.RouterGroup.POST(uri, h.h())
}

// DELETE .
func (r *Router) DELETE(uri string, h Handler) {
	r.RouterGroup.DELETE(uri, h.h())
}
