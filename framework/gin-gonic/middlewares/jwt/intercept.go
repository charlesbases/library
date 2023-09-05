package jwt

import (
	"net/http"
	"strings"
)

// ir .
type ir struct {
	// required 路径强匹配
	required []string
	// preferred 前缀匹配
	preferred []string
}

// match .
func (ir *ir) match(uri *string) bool {
	// required
	for _, val := range ir.required {
		if val == *uri {
			return true
		}
	}

	// preferred
	for _, val := range ir.preferred {
		if strings.HasPrefix(*uri, val) {
			return true
		}
	}

	return false
}

// Interceptor uri 拦截器
type Interceptor struct {
	// Enabled true: 使用 Interceptor 规则; false: 不使用 Interceptor 规则, 全部路由进行 token 校验
	Enabled bool `yaml:"enabled"`

	// Includes 进行 token 校验。优先级：高
	Includes []string `yaml:"includes"`
	includes *ir

	// Excludes 不进行 token 校验。优先级：低
	Excludes []string `yaml:"excludes"`
	excludes *ir
}

// intercept 是否进行 token 校验
func (ir *Interceptor) intercept(r *http.Request) bool {
	if !ir.Enabled {
		return true
	}

	uri := strings.Split(r.RequestURI, "?")[0]

	if ir.includes.match(&uri) {
		return true
	}

	return !ir.excludes.match(&uri)
	//
	// for _, exclude := range ir.Excludes {
	// 	// uri 全匹配
	// 	if len(exclude) == len(uri) && uri == exclude {
	// 		return false
	// 	}
	// 	// 前缀匹配
	// 	if strings.HasSuffix(exclude, "*") {
	// 		if strings.HasPrefix(uri, strings.TrimSuffix(exclude, "*")) {
	// 			return false
	// 		}
	// 	}
	// }
	//
	// return true
}
