package transport

import (
	"context"
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

type ExactMessage[T any] interface {
	proto.Message
	*T
}

var TransportModule = fx.Options(
	fx.Provide(NewEtcdTransportServer),
)

type Result[Resp any] struct {
	Client  string
	Success Resp
	Err     error
}

type EtcdTransportServer struct {
	etcdClient *etcdclient.Client
}

type Request struct {
	ID     string
	Client string
	Data   []byte
}

type Response struct {
	ID     string
	Client string
	Data   []byte
	Error  error
}

func NewEtcdTransportServer(client *etcdclient.Client) *EtcdTransportServer {
	return &EtcdTransportServer{
		etcdClient: client,
	}
}

func SendRequests[RespValue any, Resp ExactMessage[RespValue]](t *EtcdTransportServer, agents []string, msg proto.Message) ([]Result[*RespValue], error) {
	respCh := make(chan *Response, len(agents))

	for _, agent := range agents {
		go func(agentName string) {
			resp, err := t.SendRequest(agentName, msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to send request to agent")
				respCh <- &Response{
					Error: err,
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
			return nil, resp.Error
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

func SendRequest[RespValue any, Resp ExactMessage[RespValue]](t *EtcdTransportServer, client string, msg proto.Message) (*RespValue, error) {
	resp, err := t.SendRequest(client, msg)
	if err != nil {
		return nil, err
	}

	var result RespValue
	if err := proto.Unmarshal(resp.Data, Resp(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *EtcdTransportServer) SendRequest(client string, msg proto.Message) (*Response, error) {

	anyreq, err := anypb.New(msg)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rawReq, err := proto.Marshal(anyreq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req := Request{
		ID:     uuid.NewString(),
		Client: client,
		Data:   rawReq,
	}

	path := path.Join(RPCBasePath, RPCRequestPath, req.Client, req.ID)

	lease, err := t.etcdClient.Grant(context.Background(), 30)
	if err != nil {
		return nil, fmt.Errorf("failed to grant lease: %w", err)
	}

	_, err = t.etcdClient.Put(context.Background(), path, string(rawReq), clientv3.WithLease(lease.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to etcd: %w", err)
	}

	return t.waitForResponse(req)
}

func (t *EtcdTransportServer) waitForResponse(req Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	responePath := path.Join(RPCBasePath, RPCResponsePath, req.Client, req.ID)

	resp, err := t.etcdClient.Get(ctx, responePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from etcd: %w", err)
	}

	if len(resp.Kvs) > 0 {
		return &Response{
			Data: resp.Kvs[0].Value,
		}, nil
	}

	watchCh := t.etcdClient.Watch(ctx, responePath)
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
						Data: event.Kv.Value,
					}, nil
				}
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("context deadline exceeded")
		}
	}
}
