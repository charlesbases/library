package gin_gonic

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/charlesbases/logger"
	"github.com/charlesbases/logger/filewriter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/broker/kafka"
	"github.com/charlesbases/library/broker/nats"
	"github.com/charlesbases/library/codec/yaml"
	"github.com/charlesbases/library/database"
	"github.com/charlesbases/library/database/orm"
	"github.com/charlesbases/library/database/orm/driver"
	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
	"github.com/charlesbases/library/framework/gin-gonic/middlewares"
	"github.com/charlesbases/library/framework/gin-gonic/middlewares/jwt"
	"github.com/charlesbases/library/framework/gin-gonic/websocket"
	"github.com/charlesbases/library/jwtauth"
	"github.com/charlesbases/library/lifecycle"
	"github.com/charlesbases/library/redis"
	"github.com/charlesbases/library/regexp"
	"github.com/charlesbases/library/storage"
	"github.com/charlesbases/library/storage/s3"
	"github.com/charlesbases/library/watchdog"
)

// configuration .
type configuration struct {
	// Name server name
	Name string `yaml:"name"`
	// Port http port
	Port string `yaml:"port" default:":8080"`
	// Spec spec
	Spec spec `yaml:"spec"`
	// Data 服务自定义配置
	Data interface{} `yaml:"data"`
}

// spec .
type spec struct {
	// Watchdog watchdog
	Watchdog autogc `yaml:"watchdog"`
	// JWT jwt
	JWT webtoken `yaml:"jwt"`
	// Logging logging
	Logging logging `yaml:"logging"`
	// Metrics metrics
	Metrics metrics `yaml:"metrics"`
	// WebSocket websocket
	WebSocket ws `yaml:"websocket"`
	// Plugins plugins
	Plugins plugins `yaml:"plugins"`
}

// autogc .
type autogc struct {
	// Enable watchdog
	Enable bool `yaml:"enable"`
}

// logging .
type logging struct {
	OutputPath string `yaml:"outputPath"`
	MaxRolls   int    `yaml:"maxRolls"`
	Minlevel   string `yaml:"minlevel"`
	Colourful  bool   `yaml:"colourful"`
}

// metrics .
type metrics struct {
	// Enabled enabled
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// ws .
type ws struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// EnSubscription 是否启用消息订阅
	EnSubscription bool `yaml:"enSubscription"`
}

// webtoken .
type webtoken struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// Secret jwt secret
	Secret string `yaml:"secret"`
	// Expire token 过期时间。单位：秒
	Expire int `yaml:"expire"`
	// Interceptor jwt 拦截器
	Interceptor *jwt.Interceptor `yaml:"intercept"`
}

// plugins .
type plugins struct {
	// Redis redis
	Redis pluginRedis `yaml:"redis"`
	// Broker broker
	Broker pluginBroker `yaml:"broker"`
	// Storage storage
	Storage pluginStorage `yaml:"storage"`
	// Database database
	Database pluginDatabase `yaml:"database"`
}

// pluginRedis .
type pluginRedis struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// Type client or cluster
	Type string `yaml:"type"`
	// Address address for redis
	Address []string `yaml:"address"`
	// Username username
	Username string `yaml:"username"`
	// Password password
	Password string `yaml:"password"`
	// Timeout timeout
	Timeout int `yaml:"timeout" default:"3"`
	// MaxRetries 命令执行失败时的最大重试次数
	MaxRetries int `yaml:"maxRetries"`
}

// pluginBroker .
type pluginBroker struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// Type type of broker
	Type string `yaml:"type"`
	// Version kafka version
	Version string `yaml:"version"`
	// Address address
	Address string `yaml:"address"`
	// ReconnectWait default: 3s
	ReconnectWait int `yaml:"reconnectWait" default:"3"`
}

// pluginStorage .
type pluginStorage struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// Type storage.Type
	Type string `yaml:"type"`
	// Address address
	Address string `yaml:"address"`
	// AccessKey accesskey
	AccessKey string `yaml:"accessKey"`
	// SecretKey secretkey
	SecretKey string `yaml:"secretKey"`
	// Timeout timeout
	Timeout int `yaml:"timeout" default:"3"`
	// UseSSL usessl
	UseSSL bool `yaml:"useSsl"`
}

// pluginDatabase .
type pluginDatabase struct {
	// Enabled enabled
	Enabled bool `yaml:"enabled"`
	// Type database.Driver
	Type string `yaml:"type"`
	// Dsn database dsn
	Dsn string `yaml:"dsn"`
	// MaxOpenConns 最大连接数
	MaxOpenConns int `yaml:"maxOpenConns" default:"0"`
	// MaxIdleConns 连接池中最大空闲数
	MaxIdleConns int `yaml:"maxIdleConns" default:"4"`
}

// engine .
func (c *configuration) engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()
	e.Use(middlewares.Cors())
	e.Use(middlewares.Recovery())
	e.Use(middlewares.Negroni.HandlerFunc())

	// logging
	logger.SetDefault(func(o *logger.Options) {
		conf := c.Spec.Logging
		o.Colourful = conf.Colourful
		o.MinLevel = conf.Minlevel
		o.ContextHook = hfwctx.ContextHook
		o.Writer = filewriter.New(func(o *filewriter.Options) {
			o.FilePath = conf.OutputPath
			o.MaxRolls = conf.MaxRolls
		})
	})

	// jwt
	if c.Spec.JWT.Enabled {
		// init jwt
		jwtauth.Set(c.Spec.JWT.Secret, jwtauth.Expire(c.Spec.JWT.Expire))
		// use middlewares
		e.Use(jwt.New(func(j *jwt.JwtHandler) { j.Interceptor = c.Spec.JWT.Interceptor }).HandlerFunc())
	}

	// metrics
	if c.Spec.Metrics.Enabled {
		e.GET(c.Spec.Metrics.Path, gin.WrapH(promhttp.Handler()))
		middlewares.Negroni.Ignore(c.Spec.Metrics.Path)
	}

	return e
}

// redis .
func (c *configuration) redis(id string) *lifecycle.Hook {
	if !c.Spec.Plugins.Redis.Enabled {
		return nil
	}

	return &lifecycle.Hook{
		Name: "redis",
		OnStart: func(ctx context.Context) error {
			var cmdable redis.Cmdable
			switch c.Spec.Plugins.Redis.Type {
			case "client":
				cmdable = redis.RedisClient
			case "cluster":
				cmdable = redis.RedisCluster
			default:
				return errors.Errorf(`load configuration failed: unsupported values of 'spec.plugins.redis.type: "%s"'`, c.Spec.Plugins.Database.Type)
			}

			return redis.Init(id, func(o *redis.Options) {
				o.Cmdable = cmdable
				o.Addrs = c.Spec.Plugins.Redis.Address
				o.Username = c.Spec.Plugins.Redis.Username
				o.Password = c.Spec.Plugins.Redis.Password
				o.Timeout = time.Duration(c.Spec.Plugins.Redis.Timeout) * time.Second
				o.MaxRetries = c.Spec.Plugins.Redis.MaxRetries
			})
		},
		OnStop: func(ctx context.Context) error {
			if redis.Client() != nil {
				return redis.Close()
			}
			return nil
		},
	}
}

// broker .
func (c *configuration) broker(id string) *lifecycle.Hook {
	if !c.Spec.Plugins.Broker.Enabled {
		return nil
	}

	return &lifecycle.Hook{
		Name: c.Spec.Plugins.Broker.Type,
		OnStart: func(ctx context.Context) error {
			switch c.Spec.Plugins.Broker.Type {
			case "nats":
				return broker.Init(nats.NewClient(id, func(o *broker.Options) {
					o.Address = c.Spec.Plugins.Broker.Address
					o.ReconnectWait = time.Duration(c.Spec.Plugins.Broker.ReconnectWait) * time.Second
				}))
			case "kafka":
				return broker.Init(kafka.NewClient(id, func(o *broker.Options) {
					o.Version = c.Spec.Plugins.Broker.Version
					o.Address = c.Spec.Plugins.Broker.Address
					o.ReconnectWait = time.Duration(c.Spec.Plugins.Broker.ReconnectWait) * time.Second
				}))
			default:
				return errors.Errorf(`load configuration failed: unsupported values of 'spec.plugins.broker.type: "%s"'`, c.Spec.Plugins.Broker.Type)
			}
		},
		OnStop: func(ctx context.Context) error {
			if broker.C != nil {
				broker.C.Close()
			}
			return nil
		},
	}
}

// storage .
func (c *configuration) storage() *lifecycle.Hook {
	if !c.Spec.Plugins.Storage.Enabled {
		return nil
	}

	return &lifecycle.Hook{
		Name: c.Spec.Plugins.Storage.Type,
		OnStart: func(ctx context.Context) error {
			switch c.Spec.Plugins.Storage.Type {
			case "s3":
				return storage.Init(s3.NewClient(
					c.Spec.Plugins.Storage.Address,
					c.Spec.Plugins.Storage.AccessKey,
					c.Spec.Plugins.Storage.SecretKey,
					func(o *storage.Options) {
						o.Timeout = time.Duration(c.Spec.Plugins.Storage.Timeout) * time.Second
						o.UseSSL = c.Spec.Plugins.Storage.UseSSL
					}))
			default:
				return errors.Errorf(`load configuration failed: unsupported values of 'spec.plugins.storage.type: "%s"'`, c.Spec.Plugins.Storage.Type)
			}
		},
	}
}

// database .
func (c *configuration) database() *lifecycle.Hook {
	if !c.Spec.Plugins.Database.Enabled {
		return nil
	}

	return &lifecycle.Hook{
		Name: c.Spec.Plugins.Database.Type,
		OnStart: func(ctx context.Context) error {
			var dr driver.Driver
			switch c.Spec.Plugins.Database.Type {
			case "mysql":
				dr = new(driver.Mysql)
			case "postgres":
				dr = new(driver.Postgres)
			default:
				return errors.Errorf(`load configuration failed: unsupported values of 'spec.plugins.database.type: "%s"'`, c.Spec.Plugins.Database.Type)
			}

			return orm.Init(dr, func(o *database.Options) {
				o.Address = c.Spec.Plugins.Database.Dsn
				o.MaxOpenConns = c.Spec.Plugins.Database.MaxOpenConns
				o.MaxIdleConns = c.Spec.Plugins.Database.MaxIdleConns
			})
		},
	}
}

// websocket .
func (c *configuration) websocket() *lifecycle.Hook {
	if !c.Spec.WebSocket.Enabled || !c.Spec.WebSocket.EnSubscription {
		return nil
	}

	return &lifecycle.Hook{
		Name: "websocket",
		OnStart: func(ctx context.Context) error {
			return websocket.InitStation(broker.C)
		},
	}
}

// watchdog .
func (c *configuration) watchdog() *lifecycle.Hook {
	if c.Spec.Watchdog.Enable {
		var onstop = func() {}

		return &lifecycle.Hook{
			Name: "watchdog",
			OnStart: func(ctx context.Context) error {
				onstop = watchdog.Memory()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				onstop()
				return nil
			},
		}
	}

	return nil
}

// serverid .
func (c *configuration) serverid() string {
	switch model {
	case NormalModel:
		return c.Name
	case RandomModel:
		return strings.Join([]string{c.Name, uuid.NewString()}, ".")
	case HostnameModel:
		hosname, err := os.Hostname()
		if err != nil {
			logger.Fatal(err)
		}
		return strings.Join([]string{c.Name, hosname}, ".")
	case DistributionModel:
		logger.Fatal("TODO")
	default:
		logger.Fatal("unsupported model of: ", model)
	}

	return c.Name
}

// server .
func (c *configuration) server() *Server {
	srv := &Server{
		id:        c.serverid(),
		name:      c.Name,
		port:      c.Port,
		data:      c.Data,
		lifecycle: new(lifecycle.Lifecycle),
	}

	if !regexp.ServerName.MatchString(srv.name) {
		logger.Fatalf("the server name of '%s' is not allowed, must match regular of `%s`.", srv.name, regexp.ServerName.String())
	}

	// gin.Engine
	srv.engine = c.engine()

	// watchdog
	if hook := c.watchdog(); hook != nil {
		srv.lifecycle.Append(hook)
	}

	// redis
	if hook := c.redis(srv.id); hook != nil {
		srv.lifecycle.Append(hook)
	}

	// broker
	if hook := c.broker(srv.id); hook != nil {
		srv.lifecycle.Append(hook)
	}

	// storage
	if hook := c.storage(); hook != nil {
		srv.lifecycle.Append(hook)
	}

	// database
	if hook := c.database(); hook != nil {
		srv.lifecycle.Append(hook)
	}

	// websocket
	// 若启用 websocket 的 subscribe 功能，websocket.InitStation() 需要在 broker 初始化之后调用
	if hook := c.websocket(); hook != nil {
		srv.lifecycle.Append(hook)
	}

	return srv
}

// decode parse conf with Options.ConfPath
func decode() *configuration {
	var conf = new(configuration)
	if err := yaml.NewDecoder().Decode(conf); err != nil {
		logger.Fatal(err)
	}
	return conf
}
