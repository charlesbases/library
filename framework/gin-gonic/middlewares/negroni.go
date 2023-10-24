package middlewares

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
)

var (
	defaultDateFormat = "2006-01-02 15:04:05.000"
	defaultFormat     = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Key}}"
	defaulttemplate   = template.Must(template.New("negroni_parser").Parse(defaultFormat))
)

type entry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

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
			return true
		}
	}
	return false
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

			buff := new(bytes.Buffer)
			defaulttemplate.Execute(buff, &entry{
				StartTime: start.Format(defaultDateFormat),
				Status:    ctx.Writer.Status(),
				Duration:  time.Since(start),
				Hostname:  ctx.Request.Host,
				Method:    ctx.Request.Method,
				Path:      ctx.Request.URL.Path,
				Request:   ctx.Request,
			})

			logger.WithContext(c).Info(buff.String())
		} else {
			ctx.Next()
		}
	}
}
