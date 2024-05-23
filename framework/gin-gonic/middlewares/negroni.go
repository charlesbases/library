package middlewares

import (
	"fmt"
	"time"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
)

// "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"
const format = "%s | %d | %v | %s | %s %s"

// Negroni .
var Negroni = &negroni{ignores: make([]string, 0)}

// negroni .
type negroni struct {
	ignores []string
}

// allowed  是否显示请求的 info 日志
func (n *negroni) allowed(uri string) bool {
	if len(n.ignores) == 0 {
		return true
	}
	for _, val := range n.ignores {
		if uri == val {
			return false
		}
	}
	return true
}

// Ignore .
func (n *negroni) Ignore(uri ...string) {
	if len(n.ignores) == 0 {
		n.ignores = uri
	} else {
		n.ignores = append(n.ignores, uri...)
	}
}

// HandlerFunc .
func (n *negroni) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if n.allowed(ctx.Request.URL.Path) {
			start := time.Now()

			c := hfwctx.Encode(ctx)
			c.Next()

			logger.Context(c).Info(fmt.Sprintf(
				format,
				library.TimeFormat(start),
				ctx.Writer.Status(),
				time.Since(start),
				ctx.Request.RemoteAddr,
				ctx.Request.Method,
				ctx.Request.URL.Path,
			))
		} else {
			ctx.Next()
		}
	}
}
