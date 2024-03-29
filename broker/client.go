package broker

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var (
	// defaultContext default context
	defaultContext = context.Background()
	// defaultCallerSkip default of func caller skip
	defaultCallerSkip = 1
	// defaultReconnectWait 重连等待时间
	defaultReconnectWait = time.Second * 3

	// ErrNotReady .
	ErrNotReady = errors.New("connection not ready")
	// ErrInvalidAddrs .
	ErrInvalidAddrs = errors.New("invalid addrs")
)

// Client .
type Client interface {
	// Publish 消息发布
	// 若发布消息格式为 json, 则参数 'v' 为 JsonMessage.Data
	// 若发布消息格式为 proto, 则参数 'v' 为发布的完整消息体, 方法内不做额外的封装
	Publish(topic string, v interface{}, opts ...func(o *PublishOptions)) error
	// Subscribe 消息订阅
	Subscribe(topic string, handler Handler, opts ...func(o *SubscribeOptions)) error
	// Close .
	Close()
}

// Handler with Subscribe
type Handler func(event Event) error

// JsonMessage .
type JsonMessage struct {
	// ID topic 唯一标识符. uuid.NewString()
	ID string `json:"id"`
	// Producer .
	Producer string `json:"producer"`
	// CreatedAt 创建时间
	CreatedAt string `json:"created_at"`
	// Data .
	Data interface{} `json:"data"`
}

// CheckSubject .
func CheckSubject(t string) error {
	if len(strings.TrimSpace(t)) == 0 {
		return errors.New("topic cannot be empty")
	}
	if strings.Contains(t, ".") {
		return errors.New("topic cannot contain '.'")
	}
	if !utf8.ValidString(t) {
		return errors.New("topic with non UTF-8 strings are not supported")
	}
	return nil
}

// C default client
var C Client

// Init .
func Init(c Client, e error) error {
	C = c
	return e
}
