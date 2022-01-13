package registry

import (
	"errors"

	"library/registry/pb"
)

var (
	ErreEptyNode   = errors.New("require at least one node")
	ErrMissingPort = errors.New("missing port in address")
)

type Registry interface {
	Init(opts ...Option) error
	Options() *Options
	Register(*pb.Service, ...RegisterOption) error
	Deregister(*pb.Service, ...DeregisterOption) error
	GetService(string, ...ListOption) ([]*pb.Service, error)
	ListServices(...ListOption) ([]*pb.Service, error)
	String() string
}
