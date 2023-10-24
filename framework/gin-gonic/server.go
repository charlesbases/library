package gin_gonic

import (
	"os"

	"github.com/charlesbases/logger"
	"github.com/gin-gonic/gin"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/codec/json"
	"github.com/charlesbases/library/lifecycle"
)

var model = NormalModel

const (
	// NormalModel 单副本时将 name 作为服务唯一 id
	NormalModel int8 = iota
	// RandomModel 多副本但不需要连接 nats、kafka 等有状态应用时，可使用 "name.uuid" 作为服务随机唯一 id。
	// 以 nats 为例，消息订阅时的 queueName(id.subject) 和 consumerName(subject.id) 需确保唯一且不变，
	// 否则 nats 不会推送服务上次离线期间的消息。会造成消息丢失现象。
	RandomModel
	// HostnameModel 在副本 hostname 确定或规范时(StatefulSet)，"name.hostname" 可作为单副本或多副本有状态下的服务唯一 id，
	// 在 hostname 随机时(Deployment)，不可作为单副本或多副本模式下的服务唯一 id。
	HostnameModel
	// DistributionModel 使用第三方服务分发多副本的服务唯一 id
	DistributionModel
)

// Server .
type Server struct {
	id   string
	name string
	port string
	data interface{}

	lifecycle *lifecycle.Lifecycle

	engine *gin.Engine
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
func (srv *Server) Publish(topic string, v interface{}, opts ...func(o *broker.PublishOptions)) error {
	return broker.C.Publish(topic, v, opts...)
}

// Subscribe 消息异步订阅
func (srv *Server) Subscribe(topic string, handler broker.Handler, opts ...func(o *broker.SubscribeOptions)) error {
	return broker.C.Subscribe(topic, handler, opts...)
}

// Unmarshal 序列化 configuration 中的自定义配置项
func (srv *Server) Unmarshal(v interface{}) error {
	data, err := json.Marshaler.Marshal(&srv.data)
	if err != nil {
		return err
	}

	return json.Marshaler.Unmarshal(data, v)
}

// SetModel .
// NormalModel | RandomModel | HostnameModel | DistributionModel
func SetModel(m int8) {
	model = m
}

// Run .
func Run(fn func(srv *Server)) {
	// new server
	srv := decode().server()

	// do something
	fn(srv)

	// on start
	if err := srv.lifecycle.Start(); err != nil {
		srv.stop()
	}

	// on stop
	go func() {
		<-library.Shutdown()
		srv.stop()
	}()

	// run
	logger.Infof("[%s] listening and serving HTTP on %s", srv.id, srv.port)

	if err := srv.engine.Run(srv.port); err != nil {
		logger.Fatal(err)
	}
}

// stop .
func (srv *Server) stop() {
	srv.lifecycle.Stop()
	logger.Flush()

	os.Exit(1)
}
