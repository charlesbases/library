package hfwctx

import (
	"io"
	"net/http"
	"strings"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/server/webserver"
	"github.com/charlesbases/library/sonyflake"
)

// Context .
type Context struct {
	*gin.Context

	id sonyflake.ID

	log *logger.Logger
}

// ID .
func (c *Context) ID() sonyflake.ID {
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
	var id sonyflake.ID
	// parse traceid in header
	if v := ctx.GetHeader(library.HeaderTraceID); len(v) != 0 {
		id = sonyflake.ParseString(v)
	} else {
		id = sonyflake.NextID()
		// set traceid in header
		ctx.Header(library.HeaderTraceID, id.String())
	}

	// set traceid in context
	ctx.Set(library.HeaderTraceID, id.String())
	return &Context{Context: ctx, id: id, log: logger.Named(id.String(), func(o *logger.Options) { o.Skip = 1 })}
}
