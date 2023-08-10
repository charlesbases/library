package hfwctx

import (
	"io"
	"net/http"
	"strings"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/server/webserver"
)

// Context .
type Context struct {
	*gin.Context

	id  string
	log *logger.Logger
}

// ID .
func (c *Context) ID() string {
	return c.id
}

// returnJson .
func (c *Context) returnJson(v interface{}) {
	c.JSON(http.StatusOK, v)
}

// ReturnData .
func (c *Context) ReturnData(v interface{}) {
	c.returnJson(&webserver.Response{
		ID:   c.id,
		Code: webserver.StatusOK,
		Data: v,
	})
}

// ReturnError return weberror
func (c *Context) ReturnError(code webserver.Code, message ...string) {
	if len(message) != 0 {
		var mess strings.Builder
		mess.WriteString(code.Error() + ".")

		for _, val := range message {
			mess.WriteString(" ")
			mess.WriteString(val)
		}

		c.returnJson(&webserver.WebError{
			ID:      c.id,
			Code:    code,
			Message: mess.String(),
		})
	} else {
		c.returnJson(&webserver.WebError{
			ID:      c.id,
			Code:    code,
			Message: code.Error(),
		})
	}
}

// ReturnStream .
func (c *Context) ReturnStream(reader io.Reader, opts ...func(c *Context)) {
	c.Header("Content-Type", content.Stream.String())
	// c.Header("Content-Disposition", "filename")
	for _, opt := range opts {
		opt(c)
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.ReturnError(webserver.StatusServerError, err.Error())
		return
	}
}

// Decode .
func (c *Context) Decode() *gin.Context {
	return c.Context
}

// Encode .
func Encode(ctx *gin.Context) *Context {
	// parse traceid in header
	var id = ctx.GetHeader(library.HeaderTraceID)
	if len(id) == 0 {
		id = uuid.NewString()
		// set traceid in header
		ctx.Header(library.HeaderTraceID, id)
	}

	// set traceid in context
	ctx.Set(library.HeaderTraceID, id)
	return &Context{Context: ctx, id: id, log: logger.Named(id, func(o *logger.Options) { o.Skip = 1 })}
}
