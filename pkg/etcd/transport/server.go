package transport

import (
	"context"
	"errors"
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
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
)

// TransportServerModule is the server fx provider for etcd transport.
var TransportServerModule = fx.Options(
	fx.Provide(NewEtcdTransportServer),
)

// EtcdTransportServer is the server side for the etcd transport.
type EtcdTransportServer struct {
	ctx        context.Context
	cancel     context.CancelFunc
	etcdClient *etcdclient.Client
	Registry   *HandlerRegistry
	waitGroup  panichandler.WaitGroup
}

// NewEtcdTransportServer creates and returns a new etcd transport server module.
func NewEtcdTransportServer(client *etcdclient.Client) (*EtcdTransportServer, error) {
	if client == nil {
		return nil, errors.New("provided etcd client is nil")
	}
	c := &EtcdTransportServer{
		etcdClient: client,
		Registry:   NewHandlerRegistry(),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c, nil
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

// RegisterWatcher allows to register a server on the etcd transport.
func RegisterWatcher(lc fx.Lifecycle, t *EtcdTransportServer, agentName string) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			t.waitGroup.Go(func() {
				for {
					err := t.RegisterWatcher(agentName)
					if t.ctx.Err() != nil {
						log.Info().Err(t.ctx.Err()).Msg("Context canceled, stopping etcd watcher")
						return
					}
					log.Error().Err(err).Msg("etcd watch channel was canceled. Re-starting watcher")
				}
			})
			return nil
		},
		OnStop: func(_ context.Context) error {
			t.Stop()
			return nil
		},
	})
}

// RegisterWatcher register an agent on the etcd transport.
func (c *EtcdTransportServer) RegisterWatcher(agentName string) error {
	path := path.Join(RPCBasePath, RPCRequestPath, agentName)
	watchCh, err := c.etcdClient.Watch(c.ctx, path, clientv3.WithPrefix())
	if err != nil {
		log.Error().Err(err).Msg("failed to watch etcd path")
		return err
	}
	for watchResp := range watchCh {
		if watchResp.Err() != nil {
			log.Error().Err(watchResp.Err()).Msg("failed to watch etcd path")
			return watchResp.Err()
		}

		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypePut {
				id, _ := strings.CutPrefix(string(event.Kv.Key), path)
				request := Request{
					ID:     id,
					Data:   event.Kv.Value,
					Client: agentName,
				}
				go c.handleRequest(c.ctx, request)
			}
		}
	}
	return nil
}

// Stop stops the etcd transport server.
func (c *EtcdTransportServer) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}

func (c *EtcdTransportServer) handleRequest(ctx context.Context, req Request) {
	var msg anypb.Any
	err := proto.Unmarshal(req.Data, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response data")
		return
	}
	result, err := c.callHandler(ctx, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to handle request")
		return
	}
	response := Response{
		Client: req.Client,
		Data:   result,
		ID:     req.ID,
	}
	c.respond(ctx, response)
}

func (c *EtcdTransportServer) callHandler(ctx context.Context, req *anypb.Any) ([]byte, error) {
	handler, exists := c.Registry.handlers[req.MessageName()]
	if !exists {
		return nil, status.Error(
			codes.Unavailable,
			fmt.Sprintf("no handler for type %s", req.MessageName()),
		)
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling handler for request: %w", err)
	}

	serializedResp, err := proto.Marshal(resp)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return serializedResp, nil
}

func (c *EtcdTransportServer) respond(ctx context.Context, resp Response) {
	path := path.Join(RPCBasePath, RPCResponsePath, resp.Client, resp.ID)

	c.etcdClient.PutWithExpiry(path, string(resp.Data), 30)
}

// RegisterFunction register a function as a handler in the registry
// Only one function for a given Req type can be registered.
func RegisterFunction[Req, Resp proto.Message](
	t *EtcdTransportServer,
	handler func(context.Context, Req) (Resp, error),
) error {
	if handler == nil {
		return errors.New("handler cannot be nil")
	}

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
