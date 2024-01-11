package hfwctx

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/charlesbases/logger"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/charlesbases/library"
)

const _defaultConnectTimeout = 3 * time.Second

// sync.Pool of options
var pool = sync.Pool{
	New: func() interface{} {
		return &options{
			cli: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: _defaultConnectTimeout,
			},
			ctx:     context.Background(),
			body:    new(bytes.Buffer),
			headers: make(map[string]string),
		}
	},
}

// Data bytes
type Data []byte

// Unmarshal .
func (data Data) Unmarshal(pointer interface{}) error {
	if v, ok := pointer.(proto.Message); ok {
		return proto.Unmarshal(data, v)
	}
	return json.Unmarshal(data, pointer)
}

// options .
type options struct {
	cli *http.Client
	ctx context.Context

	body    *bytes.Buffer
	params  []string
	headers map[string]string
}

// option .
type option func(o *options)

// warp url
func (opts *options) warp(host string) string {
	if len(opts.params) == 0 {
		return host
	}
	return strings.Join([]string{host, strings.Join(opts.params, "&")}, "?")
}

// free .
func (opts *options) free() {
	opts.cli.Timeout = _defaultConnectTimeout

	opts.ctx = context.Background()
	opts.body.Reset()

	if cap(opts.params) > 16 {
		opts.params = nil
	} else {
		opts.params = opts.params[:0]
	}

	for key := range opts.headers {
		delete(opts.headers, key)
	}

	pool.Put(opts)
}

// newOptions .
func newOptions() *options {
	return pool.Get().(*options)
}

// Body .
func Body(data []byte) option {
	return func(o *options) {
		o.body.Write(data)
	}
}

// Param .
func Param(key string, val ...interface{}) option {
	return func(o *options) {
		if len(val) == 0 {
			o.params = append(o.params, key+"=")
		} else {
			for idx := range val {
				o.params = append(o.params, fmt.Sprintf(`%s=%v`, key, val[idx]))
			}
		}
	}
}

// Header .
func Header(key string, val string) option {
	return func(o *options) {
		o.headers[key] = val
	}
}

// Context .
func Context(ctx context.Context) option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// Timeout .
func Timeout(d time.Duration) option {
	return func(o *options) {
		o.cli.Timeout = d
	}
}

// Get .
func Get(path string, options ...option) (Data, error) {
	return newRequest(http.MethodGet, path, options...)
}

// Put .
func Put(path string, options ...option) (Data, error) {
	return newRequest(http.MethodPut, path, options...)
}

// Post .
func Post(path string, options ...option) (Data, error) {
	return newRequest(http.MethodPost, path, options...)
}

// Delete .
func Delete(path string, options ...option) (Data, error) {
	return newRequest(http.MethodDelete, path, options...)
}

// newRequest .
func newRequest(method string, host string, options ...option) (Data, error) {
	opts := newOptions()
	defer opts.free()

	for _, o := range options {
		o(opts)
	}

	req, err := http.NewRequestWithContext(opts.ctx, method, opts.warp(host), opts.body)
	if err != nil {
		return nil, errors.Errorf(`[http] %s "%s": %v`, method, strings.Split(host, "?")[0], err)
	}

	return opts.do(req)
}

// do .
func (opts *options) do(req *http.Request) (Data, error) {
	for key, val := range opts.headers {
		req.Header.Set(key, val)
	}

	start := time.Now()

	var retry int8
conn:
	rsp, err := opts.cli.Do(req)
	if err != nil {
		// 如果连接超时，则重试
		if uerr, ok := err.(*url.Error); ok && uerr.Timeout() && retry < 3 {
			retry++
			goto conn
		}
		return nil, errors.Errorf("[http] %v", err)
	}
	defer rsp.Body.Close()

	logger.CallerSkip(3).WithContext(req.Context()).Debugf(
		"[http] %s | %d | %v | %s | %s %s",
		library.TimeFormat(start), rsp.StatusCode, time.Since(start), req.URL.Host, req.Method, req.URL.Path,
	)

	switch rsp.StatusCode {
	case http.StatusOK:
		var data bytes.Buffer
		if _, err := io.Copy(&data, rsp.Body); err != nil {
			return nil, errors.Errorf(`[http] %s "%s%s": %v`, req.Method, req.URL.Host, req.URL.Path, err)
		} else {
			return data.Bytes(), nil
		}
	default:
		return nil, errors.Errorf(`[http] %s "%s%s": %s`, req.Method, req.URL.Host, req.URL.Path, rsp.Status)
	}
}
