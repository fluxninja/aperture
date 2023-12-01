package transport

import (
	"context"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

const (
	leaseDuration       = 30
	serverWatchDuration = 60
)

// ExactMessage is like proto.Message, but known to point to T.
type ExactMessage[T any] interface {
	proto.Message
	*T
}

// TransportClientModule is the client fx provider for etcd transport.
var TransportClientModule = fx.Options(
	fx.Provide(NewEtcdTransportClient),
)

// Result is a result from one agent.
type Result[Resp any] struct {
	Client  string
	Success Resp
	Err     error
}

// EtcdTransportClient is the client side of the etcd transport.
type EtcdTransportClient struct {
	etcdClient *etcdclient.Client
}

// Request is the raw request on the etcd transport.
type Request struct {
	ID     string
	Client string
	Data   []byte
}

// Response is the raw response on the etcd transport.
type Response struct {
	ID     string
	Client string
	Data   []byte
	Error  error
}

// NewEtcdTransportClient creates a new client on the etcd transport.
func NewEtcdTransportClient(client *etcdclient.Client) (*EtcdTransportClient, error) {
	if client == nil {
		return nil, errors.New("provided etcd client is nil")
	}
	return &EtcdTransportClient{
		etcdClient: client,
	}, nil
}

// SendRequests allows consumers of the etcd transport to send requests to agents.
func SendRequests[RespValue any, Resp ExactMessage[RespValue]](ctx context.Context, t *EtcdTransportClient, agents []string, msg proto.Message) ([]Result[*RespValue], error) {
	respCh := make(chan *Response, len(agents))

	for _, agent := range agents {
		go func(agentName string) {
			resp, err := t.SendRequest(ctx, agentName, msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to send request to agent")
				respCh <- &Response{
					Client: agentName,
					Error:  err,
				}
			} else {
				respCh <- resp
			}
		}(agent)
	}

	resps := make([]Result[*RespValue], 0, len(agents))
	for i := 0; i < len(agents); i++ {
		resp := <-respCh
		if resp.Error != nil {
			resps = append(resps, Result[*RespValue]{
				Err: resp.Error,
			})
			continue
		}
		var result RespValue
		if err := proto.Unmarshal(resp.Data, Resp(&result)); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		resps = append(resps, Result[*RespValue]{
			Client:  resp.Client,
			Success: &result,
		})
	}

	return resps, nil
}

// SendRequest allows consumers of the etcd transport to send single request to agents.
func SendRequest[RespValue any, Resp ExactMessage[RespValue]](ctx context.Context, t *EtcdTransportClient, agentName string, msg proto.Message) (*RespValue, error) {
	resp, err := t.SendRequest(ctx, agentName, msg)
	if err != nil {
		return nil, err
	}

	var result RespValue
	if err := proto.Unmarshal(resp.Data, Resp(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

// SendRequest sends a request to etcd, supposed to be consumed by an agent.
func (t *EtcdTransportClient) SendRequest(ctx context.Context, agentName string, msg proto.Message) (*Response, error) {
	anyreq, err := anypb.New(msg)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rawReq, err := proto.Marshal(anyreq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request %v: %w", msg, err)
	}

	req := Request{
		ID:     uuid.NewString(),
		Client: agentName,
		Data:   rawReq,
	}

	path := path.Join(RPCBasePath, RPCRequestPath, req.Client, req.ID)

	t.etcdClient.PutWithExpiry(path, string(rawReq), leaseDuration)

	return t.waitForResponse(ctx, req)
}

func (t *EtcdTransportClient) waitForResponse(ctx context.Context, req Request) (*Response, error) {
	watchCtx, cancel := context.WithTimeout(ctx, serverWatchDuration*time.Second)
	defer cancel()

	responsePath := path.Join(RPCBasePath, RPCResponsePath, req.Client, req.ID)

	watchCh, err := t.etcdClient.Watch(watchCtx, responsePath)
	if err != nil {
		return nil, fmt.Errorf("failed to watch etcd path: %w", err)
	}
	for {
		select {
		case watchResp, ok := <-watchCh:
			if !ok {
				return nil, fmt.Errorf("watch channel closed")
			}
			if watchResp.Err() != nil {
				return nil, fmt.Errorf("failed to watch etcd path: %w", watchResp.Err())
			}

			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypePut {
					return &Response{
						Client: req.Client,
						Data:   event.Kv.Value,
					}, nil
				}
			}
		case <-watchCtx.Done():
			return nil, fmt.Errorf("context deadline exceeded")
		}
	}
}
