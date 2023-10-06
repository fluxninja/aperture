package transport

import (
	"context"
	"fmt"
	"path"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// TransportClientModule is the client fx provider for etcd transport
var TransportClientModule = fx.Options(
	fx.Provide(NewEtcdTransportClient),
)

// EtcdTransportClient is the client side for the etcd transport
type EtcdTransportClient struct {
	etcdClient *etcdclient.Client
	Registry   *HandlerRegistry
}

// NewEtcdTransportClient creates and returns a new etcd transport client module
func NewEtcdTransportClient(client *etcdclient.Client) *EtcdTransportClient {
	return &EtcdTransportClient{
		etcdClient: client,
		Registry:   NewHandlerRegistry(),
	}
}

// HandlerRegistry allow registering handlers and can start a dispatcher.
//
// This is intended to be used at at fx provide/invoke stage.
type HandlerRegistry struct {
	handlers map[protoreflect.FullName]untypedHandler
}

type untypedHandler func(context.Context, *anypb.Any) (proto.Message, error)

// NewHandlerRegistry creates a new HandlerRegistry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: map[protoreflect.FullName]untypedHandler{},
	}
}

// RegisterWatcher allows to register a client on the etcd transport
func RegisterWatcher(lc fx.Lifecycle, t *EtcdTransportClient, agentName string) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go t.RegisterWatcher(agentName)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

// RegisterWatcher register an agent on the etcd transport client
func (c *EtcdTransportClient) RegisterWatcher(agentName string) {
	path := path.Join(RPCBasePath, RPCRequestPath, agentName)
	watchCh := c.etcdClient.Watch(context.Background(), path, clientv3.WithPrefix())
	for watchResp := range watchCh {
		if watchResp.Err() != nil {
			log.Error().Err(watchResp.Err()).Msg("failed to watch etcd path")
			continue
		}

		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypePut {
				id, _ := strings.CutPrefix(string(event.Kv.Key), path)
				request := Request{
					ID:     id,
					Data:   event.Kv.Value,
					Client: agentName,
				}
				go c.handleRequest(context.Background(), request)
			}
		}
	}
}

func (c *EtcdTransportClient) handleRequest(ctx context.Context, req Request) {
	var msg anypb.Any
	err := proto.Unmarshal(req.Data, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response data")
		return
	}
	result, _ := c.callHandler(ctx, &msg)
	response := Response{
		Client: req.Client,
		Data:   result,
		ID:     req.ID,
	}
	c.respond(ctx, response)
}

func (c *EtcdTransportClient) callHandler(ctx context.Context, req *anypb.Any) ([]byte, error) {
	handler, exists := c.Registry.handlers[req.MessageName()]
	if !exists {
		return nil, status.Error(
			codes.Unavailable,
			fmt.Sprintf("no handler for type %s", req.MessageName()),
		)
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	serializedResp, err := proto.Marshal(resp)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return serializedResp, nil
}

func (c *EtcdTransportClient) respond(ctx context.Context, resp Response) {

	path := path.Join(RPCBasePath, RPCResponsePath, resp.Client, resp.ID)

	lease, err := c.etcdClient.Grant(context.Background(), 30)
	if err != nil {
		log.Error().Err(err).Msg("failed to grant lease")
		return
	}

	_, err = c.etcdClient.Put(ctx, path, string(resp.Data), clientv3.WithLease(lease.ID))
	if err != nil {
		log.Error().Err(err).Msg("failed to write response to etcd")
	}
}

// RegisterFunction register a function as a handler in the registry
// Only one function for a given Req type can be registered.
func RegisterFunction[Req, Resp proto.Message](
	t *EtcdTransportClient,
	handler func(context.Context, Req) (Resp, error),
) error {
	var req Req // used only to pull out the message name from descriptor
	name := req.ProtoReflect().Descriptor().FullName()
	if _, exists := t.Registry.handlers[name]; exists {
		log.Error().Err(fmt.Errorf("handler for %s already registered", name))
		return fmt.Errorf("handler for %s already registered", name)
	}

	t.Registry.handlers[name] = func(ctx context.Context, anyreq *anypb.Any) (proto.Message, error) {
		var nilReq Req
		req := nilReq.ProtoReflect().New().Interface().(Req)
		if err := anyreq.UnmarshalTo(req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return handler(ctx, req)
	}

	return nil
}
