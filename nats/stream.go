package nats

import (
	"github.com/nats-io/nats.go"
)

const (
	defaultStreamName    = "Stream"
	defaultStreamSubject = "*"
)

type streamOption func(sc *nats.StreamConfig)

// newStreamConfig .
func newStreamConfig(opts ...streamOption) *nats.StreamConfig {
	streamConfig := defaultStreamConfig()
	for _, opt := range opts {
		opt(streamConfig)
	}
	return streamConfig
}

// defaultStreamConfig .
func defaultStreamConfig() *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:       defaultStreamName,
		Subjects:   []string{defaultStreamSubject},
		Retention:  nats.WorkQueuePolicy,
		MaxMsgs:    1 << 10,
		MaxBytes:   1 << 30, // 1 GiB
		Discard:    nats.DiscardOld,
		MaxMsgSize: 1 << 20, // 1 MiB
		Storage:    nats.FileStorage,
	}
}

// withStreamName .
func withStreamName(name string) streamOption {
	return func(sc *nats.StreamConfig) {
		sc.Name = name
	}
}

// withStreamSubjects .
func withStreamSubjects(topics ...string) streamOption {
	return func(sc *nats.StreamConfig) {
		sc.Subjects = topics
	}
}
