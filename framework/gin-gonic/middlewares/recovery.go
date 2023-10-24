package middlewares

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"
)

// Recovery .
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1<<13)
				stack = stack[:runtime.Stack(stack, false)]
				infor := &p{err: err, stack: stack, request: ctx.Request}

				logger.WithContext(ctx).Errorf("\nPanic:   [%s]\nRequest: [%s]\nStack:   [%s]",
					infor.err, infor.RequestDesc(), infor.StackAsString(),
				)
			}
		}()

		ctx.Next()
	}
}

type p struct {
	err     interface{}
	stack   []byte
	request *http.Request
}

// StackAsString .
func (p *p) StackAsString() string {
	return string(p.stack)
}

// RequestDesc .
func (p *p) RequestDesc() string {
	if p.request == nil {
		return "request is nil"
	}

	var query string
	if len(p.request.URL.RawQuery) != 0 {
		query = "?" + p.request.URL.RawQuery
	}
	return fmt.Sprintf("%s %s%s", p.request.Method, p.request.URL.Path, query)
}
