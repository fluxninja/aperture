package rpc

import (
	"context"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	rpcv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/rpc/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcd "github.com/fluxninja/aperture/v2/pkg/etcd/writer"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
)

// EtcdClient is a client part of "etcd" transport
type EtcdClient struct {
	clientName        string
	handlers          *HandlerRegistry
	dispatcher        *Dispatcher
	coordinatorClient rpcv1.CoordinatorClient
	etcdWatcher       notifiers.Watcher
	etcdWriter        etcd.Writer
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

// RegisterEtcdClient is an FX helper to connect new EtcdClient to a given addr.
func RegisterEtcdClient(
	clientName string,
	lc fx.Lifecycle,
	handlers *HandlerRegistry,
	etcdWatcher notifiers.Watcher,
	etcdWriter etcd.Writer,
	addr string,
) {
	var client *EtcdClient
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			client = NewEtcdClient(
				clientName,
				handlers,
				etcdWatcher,
				etcdWriter,
			)
			helloMsg := &rpcv1.ClientToServer{
				Msg: &rpcv1.ClientToServer_Hello_{
					Hello: &rpcv1.ClientToServer_Hello{
						Name:   client.clientName,
						NextId: 0,
					},
				},
			}

			marshalledHello, err := protojson.Marshal(helloMsg)
			if err != nil {
				return err
			}

			client.etcdWriter.Write(path.Join(paths.RPCRegistrationPathPrefix, client.clientName), marshalledHello)

			client.dispatcher = client.handlers.StartDispatcher()
			defer client.dispatcher.Stop()

			client.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			client.Stop()
			return nil
		},
	})
}

// NewEtcdClient creates a new StreamClient.
func NewEtcdClient(
	name string,
	handlers *HandlerRegistry,
	etcdWatcher notifiers.Watcher,
	etcdWriter etcd.Writer,
) *EtcdClient {
	return &EtcdClient{
		clientName:  name,
		handlers:    handlers,
		etcdWatcher: etcdWatcher,
		etcdWriter:  etcdWriter,
	}
}

// Start starts the StreamClient.
func (c *EtcdClient) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.wg.Add(1)
	go c.runClient()
}

// Stop stops the StreamClient.
func (c *EtcdClient) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *EtcdClient) runClient() {
	defer c.wg.Done()

	for {
		c.runClientIteration()

		select {
		case <-c.ctx.Done():
			return
		default:
		}

		select {
		case <-time.After(reconnectDelay):
			log.Autosample().Info().Msg("stream client: reconnecting")
			continue
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *EtcdClient) runClientIteration() {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	c.dispatcher = c.handlers.StartDispatcher()
	defer c.dispatcher.Stop()

	keyPrefixNotifier, err := notifiers.NewUnmarshalPrefixNotifier(
		paths.RPCRequestsPath(c.clientName),
		c.etcdWatchCallback,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
	)
	if err != nil {
		return
	}

	if err := c.etcdWatcher.AddPrefixNotifier(keyPrefixNotifier); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case response := <-c.dispatcher.Chan():
			c.etcdWriter.Write(path.Join(paths.RPCResponsesPath(c.clientName), strconv.FormatUint(uint64(response.Id), 10)), response.GetPayload())
		}
	}
}

func (c *EtcdClient) etcdWatchCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	if event.Type == notifiers.Remove {
		return
	}

	var msg rpcv1.ServerToClient
	err := unmarshaller.Unmarshal(&msg)
	if err != nil {
		return
	}

	c.etcdWriter.Delete(string(event.Key))

	c.dispatcher.ProcessCommand(&msg)
}

// EtcdServer is a server part of "etcd" transport
type EtcdServer struct {
	clients     *Clients
	etcdWatcher notifiers.Watcher
	etcdWriter  etcd.Writer
}

// RegisterEtcdServer is an FX helper to connect new EtcdClient to a given addr.
func RegisterEtcdServer(
	lc fx.Lifecycle,
	etcdWatcher notifiers.Watcher,
	etcdClient *etcdclient.Client,
	clients *Clients,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			server := &EtcdServer{
				clients:     clients,
				etcdWatcher: etcdWatcher,
			}
			if err := server.Start(); err != nil {
				log.Bug().Err(err)
			}

			server.etcdWriter = *etcd.NewWriter(&etcdClient.KVWrapper)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

}

// NewEtcdServer creates new EtcdServer plugging into given Clients as a transport.
func NewEtcdServer(clients *Clients, etcdWatcher notifiers.Watcher, etcdWriter etcd.Writer) *EtcdServer {
	return &EtcdServer{
		clients:     clients,
		etcdWatcher: etcdWatcher,
		etcdWriter:  etcdWriter,
	}
}

func (s *EtcdServer) etcdPrefixWatcherCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	if event.Type == notifiers.Remove {
		return
	}

	var helloMsg *rpcv1.ClientToServer
	err := unmarshaller.Unmarshal(&helloMsg)
	if err != nil {
		return
	}

	hello, ok := helloMsg.Msg.(*rpcv1.ClientToServer_Hello_)
	if !ok {
		return
	}

	requests, responses := s.clients.Join(hello.Hello.Name, hello.Hello.NextId)
	defer close(responses)

	callback := func(event notifiers.Event, unmarshaller config.Unmarshaller) {
		if event.Type == notifiers.Remove {
			return
		}

		var msg rpcv1.ClientToServer
		err = unmarshaller.Unmarshal(&msg)
		if err != nil {
			return
		}

		clientDisconnected := make(chan error, 1)

		switch msg := msg.Msg.(type) {
		case *rpcv1.ClientToServer_Response:
			responses <- msg.Response
		case *rpcv1.ClientToServer_Hello_:
			clientDisconnected <- status.Error(codes.InvalidArgument, "unexpected hello")
			return
		default:
			log.Bug().Msg("unknown client to server message")
		}

		for {
			select {
			case serverToClient, ok := <-requests:
				if !ok {
					return
				}

				marshalledReq, err := proto.Marshal(serverToClient)
				if err != nil {
					log.Bug().Err(err)
					return
				}

				s.etcdWriter.Write(path.Join(paths.RPCRequestsPath(hello.Hello.Name), strconv.FormatUint(serverToClient.GetRequest().Id, 10)), marshalledReq)

			case err := <-clientDisconnected:
				log.Bug().Err(err)
				return
			}
		}
	}

	keyNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(paths.RPCResponsesPath(hello.Hello.Name)),
		unmarshaller,
		callback,
	)
	if err != nil {
		return
	}

	if err := s.etcdWatcher.AddKeyNotifier(keyNotifier); err != nil {
		return
	}
}

func (s *EtcdServer) Start() error {
	notifier, err := notifiers.NewUnmarshalPrefixNotifier(
		paths.RPCRegistrationPathPrefix,
		s.etcdPrefixWatcherCallback,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller,
	)
	if err != nil {
		return err
	}

	s.etcdWatcher.AddPrefixNotifier(notifier)

	return nil
}
