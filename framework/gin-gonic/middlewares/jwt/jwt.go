package jwt

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library/framework/gin-gonic/webserver"
	"github.com/charlesbases/library/jwtauth"
)

const defaultTokenHeader = "Authorization"

// defaultJwtAuthHandler 默认鉴权方式。只验证 token 有效性，不进行 api 鉴权
var defaultJwtAuthHandler func(ctx *gin.Context) bool = func(ctx *gin.Context) bool {
	if _, err := UserClaims(ctx); err != nil {
		ctx.JSON(http.StatusUnauthorized, webserver.StatusUnauthorized.WebError(err))
		return false
	}
	return true
}

// JwtHandler .
type JwtHandler struct {
	// 拦截器
	Interceptor *Interceptor
	// 解析器
	Handler func(ctx *gin.Context) bool
}

// defaultHandler .
func defaultHandler() *JwtHandler {
	return &JwtHandler{
		Handler: defaultJwtAuthHandler,
	}
}

// New .
func New(opts ...func(j *JwtHandler)) *JwtHandler {
	var jwt = defaultHandler()
	for _, opt := range opts {
		opt(jwt)
	}

	// interceptor
	if jwt.Interceptor != nil && jwt.Interceptor.Enabled {
		// includes
		{
			ir := &ir{
				required:  make([]string, 0, len(jwt.Interceptor.Includes)),
				preferred: make([]string, 0, len(jwt.Interceptor.Includes)),
			}
			for _, uri := range jwt.Interceptor.Includes {
				if strings.HasSuffix(uri, "*") {
					ir.preferred = append(ir.preferred, strings.TrimSuffix(uri, "*"))
				} else {
					ir.required = append(ir.required, uri)
				}
			}
			jwt.Interceptor.includes = ir
		}

		// excludes
		{
			ir := &ir{
				required:  make([]string, 0, len(jwt.Interceptor.Excludes)),
				preferred: make([]string, 0, len(jwt.Interceptor.Excludes)),
			}
			for _, uri := range jwt.Interceptor.Excludes {
				if strings.HasSuffix(uri, "*") {
					ir.preferred = append(ir.preferred, strings.TrimSuffix(uri, "*"))
				} else {
					ir.required = append(ir.required, uri)
				}
			}
			jwt.Interceptor.excludes = ir
		}
	}
	return jwt
}

// HandlerFunc .
func (j *JwtHandler) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if j.Interceptor != nil && j.Interceptor.intercept(ctx.Request) {
			if !j.Handler(ctx) {
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

// UserClaims .
func UserClaims(ctx *gin.Context) (*jwtauth.UserClaims, error) {
	token := ctx.GetHeader(defaultTokenHeader)
	if len(token) == 0 {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	return jwtauth.Encode(token)
}
