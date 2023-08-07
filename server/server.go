package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/codec/json"
	"github.com/charlesbases/library/lifecycle"
	"github.com/charlesbases/library/storage"
)

// Server .
type Server struct {
	id   string
	name string
	uuid uuid.UUID

	ctx       context.Context
	lifecycle *lifecycle.Lifecycle

	engine *gin.Engine

	broker  broker.Client
	storage storage.Client
}

// Run .
func Run(fn func(srv *Server)) {
	// new server
	srv := parseconf().server()

	// do something
	fn(srv)

	// on start
	if err := srv.lifecycle.Start(srv.ctx); err != nil {
		os.Exit(1)
	}

	// on stop
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		select {
		case <-c:
			srv.lifecycle.Stop(srv.ctx)
			close(c)
		}
		os.Exit(1)
	}()

	// run
	logger.Infof("listening and serving HTTP on %s", conf.Port)

	if err := srv.engine.Run(conf.Port); err != nil {
		logger.Fatal(err)
	}
}

// Use .
func (srv *Server) Use(middleware ...gin.HandlerFunc) {
	srv.engine.Use(middleware...)
}

// Lifecycle .
func (srv *Server) Lifecycle(hooks ...*lifecycle.Hook) {
	srv.lifecycle.Append(hooks...)
}

// RegisterRouter 注册路由
func (srv *Server) RegisterRouter(fn func(r *Router)) {
	fn(&Router{&srv.engine.RouterGroup})
}

// RegisterRouterGroup 注册路由组
func (srv *Server) RegisterRouterGroup(uri string, fn func(r *Router), handlers ...gin.HandlerFunc) {
	fn(&Router{srv.engine.Group(uri, handlers...)})
}

// Publish 消息异步发布
func (srv *Server) Publish(topic string, v interface{}, opts ...broker.PublishOption) {
	if srv.broker != nil {
		srv.broker.Publish(topic, v, opts...)
	} else {
		logger.Warn("publish failed. no broker!")
	}
}

// Subscribe 消息订阅
func (srv *Server) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) {
	if srv.broker != nil {
		srv.broker.Subscribe(topic, handler, opts...)
	} else {
		logger.Warn("subscribe failed. no broker!")
	}
}

// Unmarshal 序列化 configuration 中的自定义配置项
func (srv *Server) Unmarshal(v interface{}) error {
	data, err := json.Marshaler.Marshal(&conf.Data)
	if err != nil {
		return err
	}

	return json.Marshaler.Unmarshal(data, v)
}
