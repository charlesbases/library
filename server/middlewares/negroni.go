package middlewares

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/charlesbases/library/server/hfwctx"
	"github.com/gin-gonic/gin"
)

var (
	defaultDateFormat = "2006-01-02 15:04:05.000"
	defaultFormat     = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"
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
func Negroni() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		c.Info(buff.String())
	}
}
