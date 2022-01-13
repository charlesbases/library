package etcd

import (
	"context"
	"crypto/tls"
	"net"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"library/registry"
	"library/registry/pb"

	"github.com/golang/protobuf/proto"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/clientv3"
)

const prefix = "/registry/"

type etcdRegistry struct {
	client  *clientv3.Client
	options *registry.Options

	sync.RWMutex
	leases map[string]int64
}

// NewRegistry returns an initialized etcd registry .
func NewRegistry(opts ...registry.Option) registry.Registry {
	e := &etcdRegistry{
		leases: make(map[string]int64, 0),
	}
	e.configure(opts...)
	return e
}

// configure .
func (e *etcdRegistry) configure(opts ...registry.Option) error {
	for _, opt := range opts {
		opt(e.options)
	}

	etcdConfig := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	{
		if e.options.Timeout == 0 {
			e.options.Timeout = 5 * time.Second
		}
		etcdConfig.DialTimeout = e.options.Timeout

		if e.options.Secure || e.options.TLSConfig != nil {
			tlsConfig := e.options.TLSConfig
			if tlsConfig == nil {
				tlsConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			}
			etcdConfig.TLS = tlsConfig
		}

		var endpoints = make([]string, 0, len(e.options.Addrs))
		for _, address := range e.options.Addrs {
			if len(address) == 0 {
				continue
			}
			addr, port, err := net.SplitHostPort(address)
			if ae, ok := err.(*net.AddrError); ok && ae.Err == registry.ErrMissingPort.Error() {
				port = "2379"
				addr = address
				endpoints = append(endpoints, net.JoinHostPort(addr, port))
			} else if err == nil {
				endpoints = append(endpoints, net.JoinHostPort(addr, port))
			}
		}
		// if we got addrs then we'll update
		if len(endpoints) > 0 {
			etcdConfig.Endpoints = endpoints
		}
		// check if the endpoints have https://
		if etcdConfig.TLS != nil {
			for i, ep := range etcdConfig.Endpoints {
				if !strings.HasPrefix(ep, "https://") {
					etcdConfig.Endpoints[i] = "https://" + ep
				}
			}
		}
	}

	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return err
	}
	if e.client != nil {
		e.client.Close()
	}
	// setup new client
	e.client = client
	return nil
}

func (e *etcdRegistry) Init(opts ...registry.Option) error {
	return e.configure(opts...)
}

func (e *etcdRegistry) Options() *registry.Options {
	return e.options
}

// getService .
func (e *etcdRegistry) getService(prefix string, options *registry.ListOptions) ([]*pb.Service, error) {
	ctx, cancel := context.WithTimeout(options.Context, e.options.Timeout)
	defer cancel()

	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	if resp.Kvs == nil || len(resp.Kvs) == 0 {
		return nil, nil
	}

	var versions = make(map[string]*pb.Service, 0)
	for _, n := range resp.Kvs {
		sn := decodeService(n.Value)
		if sn != nil {
			if s, find := versions[serviceVersion(sn.Name, sn.Version)]; find {
				s.Nodes = append(s.Nodes, sn.Nodes...)
			} else {
				versions[serviceVersion(sn.Name, sn.Version)] = sn
			}
		}
	}

	var srvs = make([]*pb.Service, 0, len(versions))
	for _, srv := range versions {
		srvs = append(srvs, srv)
	}

	sort.Slice(srvs, func(i, j int) bool { return srvs[i].Name < srvs[j].Name })
	return srvs, nil
}

func (e *etcdRegistry) GetService(name string, opts ...registry.ListOption) ([]*pb.Service, error) {
	var options = new(registry.ListOptions)
	for _, opt := range opts {
		opt(options)
	}

	return e.getService(servicePath(name), options)
}

func (e *etcdRegistry) ListServices(opts ...registry.ListOption) ([]*pb.Service, error) {
	var options = new(registry.ListOptions)
	for _, opt := range opts {
		opt(options)
	}

	return e.getService(prefix, options)
}

func (e *etcdRegistry) Register(service *pb.Service, opts ...registry.RegisterOption) error {
	if service.GetNodes() == nil || len(service.GetNodes()) == 0 {
		return registry.ErreEptyNode
	}

	var options = &registry.RegisterOptions{Context: context.Background()}
	for _, opt := range opts {
		opt(options)
	}

	var rErr error
	// register each node individually
	for _, node := range service.GetNodes() {
		if err := e.registerNode(service, node, options); err != nil {
			rErr = err
		}
	}
	return rErr
}

func (e *etcdRegistry) registerNode(s *pb.Service, n *pb.Node, options *registry.RegisterOptions) error {
	e.Lock()
	leaseID, found := e.leases[nodeName(s.Name, n.ID)]
	e.Unlock()

	// lease exists
	if found && leaseID > 0 {
		_, err := e.client.KeepAliveOnce(options.Context, clientv3.LeaseID(leaseID))
		if err == nil || err != rpctypes.ErrLeaseNotFound {
			return err
		}

		e.Lock()
		delete(e.leases, nodeName(s.Name, n.ID))
		e.Unlock()
	}

	// Register
	service := &pb.Service{
		Name:     s.GetName(),
		Version:  s.GetVersion(),
		Nodes:    []*pb.Node{n},
		Metadata: s.GetMetadata(),
	}

	ctx, cancel := context.WithTimeout(options.Context, e.options.Timeout)
	defer cancel()

	if options.TTL.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lcr, err := e.client.Create(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
		leaseID = lcr.ID
	} else {
		leaseID = -1
	}

	if _, err := e.client.Put(options.Context, nodePath(s.GetName(), n.GetID()), encodeSrevice(service)); err != nil {
		return err
	}

	e.Lock()
	e.leases[nodeName(s.GetName(), n.GetID())] = leaseID
	e.Unlock()

	return nil
}

func (e *etcdRegistry) Deregister(service *pb.Service, opts ...registry.DeregisterOption) error {
	if service.Nodes == nil || len(service.Nodes) == 0 {
		return registry.ErreEptyNode
	}

	var options = &registry.DeregisterOptions{Context: context.Background()}
	for _, opt := range opts {
		opt(options)
	}

	for _, node := range service.GetNodes() {
		ctx, cancel := context.WithTimeout(options.Context, e.options.Timeout)
		defer cancel()

		if _, err := e.client.Delete(ctx, nodePath(service.Name, node.ID)); err != nil {
			return err
		}

		e.Lock()
		delete(e.leases, nodeName(service.Name, node.ID))
		e.Unlock()
	}
	return nil
}

func (e *etcdRegistry) String() string {
	return "etcd"
}

func encodeSrevice(s *pb.Service) string {
	b, _ := proto.Marshal(s)
	return string(b)
}

func decodeService(b []byte) *pb.Service {
	var s *pb.Service
	proto.Unmarshal(b, s)
	return s
}

func nodeName(s, n string) string {
	return s + "." + n
}

func nodePath(s, n string) string {
	return path.Join(
		prefix,
		strings.ReplaceAll(s, "/", "-"),
		strings.ReplaceAll(n, "/", "-"),
	)
}

func servicePath(s string) string {
	return path.Join(
		prefix,
		strings.ReplaceAll(s, "/", "-"),
	) + "/"
}

func serviceVersion(s, v string) string {
	return s + "." + v
}
