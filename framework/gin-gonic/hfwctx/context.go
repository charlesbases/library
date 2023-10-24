package hfwctx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/framework/gin-gonic/webserver"
)

// headerTraceID traceid in context and header
const headerTraceID = "X-Trace-ID"

// ID trace id
type ID string

// String .
func (id ID) String() string {
	return string(id)
}

// NewID .
func NewID() ID {
	return ID(uuid.NewString())
}

// Context .
type Context struct {
	*gin.Context

	id ID
}

// ID .
func (c *Context) ID() ID {
	return c.id
}

// returnJson .
func (c *Context) returnJson(v interface{}) {
	c.JSON(http.StatusOK, v)
}

// Successful .
func (c *Context) Successful() {
	c.returnJson(&webserver.Response{
		ID:   c.id.String(),
		Code: webserver.StatusOK,
	})
}

// ReturnData .
func (c *Context) ReturnData(v interface{}) {
	c.returnJson(&webserver.Response{
		ID:   c.id.String(),
		Code: webserver.StatusOK,
		Data: v,
	})
}

// ReturnError return weberror
func (c *Context) ReturnError(code webserver.Code, v ...interface{}) {
	if len(v) != 0 {
		var mess strings.Builder
		mess.WriteString(code.Error() + ".")

		for _, err := range v {
			mess.WriteString(" ")
			mess.WriteString(fmt.Sprintf(`%v`, err))
		}

		c.returnJson(&webserver.WebError{
			ID:      c.id.String(),
			Code:    code,
			Message: mess.String(),
		})
	} else {
		c.returnJson(&webserver.WebError{
			ID:      c.id.String(),
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
		c.ReturnError(webserver.StatusServerError, err)
		return
	}
}

// Decode .
func (c *Context) Decode() *gin.Context {
	return c.Context
}

// Encode .
func Encode(ctx *gin.Context) *Context {
	// in context
	if id, exists := ctx.Get(headerTraceID); exists {
		return &Context{Context: ctx, id: id.(ID)}
	}

	// in header
	if id := ctx.GetHeader(headerTraceID); len(id) != 0 {
		ctx.Set(headerTraceID, id)
		return &Context{Context: ctx, id: ID(id)}
	}

	id := NewID()
	ctx.Header(headerTraceID, id.String())
	ctx.Set(headerTraceID, id)
	return &Context{Context: ctx, id: ID(id)}
}

// ContextHook .
func ContextHook(ctx context.Context) func(l *logger.Logger) *logger.Logger {
	return func(l *logger.Logger) *logger.Logger {
		if ctx == context.Background() {
			return l
		}
		if val := ctx.Value(headerTraceID); val != nil {
			return l.Named(val.(string))
		}
		return l
	}
}
