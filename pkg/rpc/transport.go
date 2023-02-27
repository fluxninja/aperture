package rpc

import (
	"context"
	"io"
	"sync"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpcv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/rpc/v1"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
)

const reconnectDelay = 10 * time.Second

// StreamClient is a client part of "stream" transport, which uses
// bidirectional grpc stream methods.
type StreamClient struct {
	clientName        string
	handlers          *HandlerRegistry
	coordinatorClient rpcv1.CoordinatorClient
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

// RegisterStreamClient is an FX helper to connect new StreamClient to a given addr.
func RegisterStreamClient(
	clientName string,
	lc fx.Lifecycle,
	handlers *HandlerRegistry,
	connWrapper grpcclient.ClientConnectionWrapper,
	addr string,
) {
	var client *StreamClient
	var conn *grpc.ClientConn
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			if conn, err = connWrapper.Dial(ctx, addr); err != nil {
				return err
			}

			client = NewStreamClient(
				clientName,
				handlers,
				rpcv1.NewCoordinatorClient(conn),
			)
			client.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			client.Stop()
			conn.Close()
			return nil
		},
	})
}

// NewStreamClient creates a new StreamClient.
func NewStreamClient(
	name string,
	handlers *HandlerRegistry,
	client rpcv1.CoordinatorClient,
) *StreamClient {
	return &StreamClient{
		clientName:        name,
		handlers:          handlers,
		coordinatorClient: client,
	}
}

// Start starts the StreamClient.
func (c *StreamClient) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.wg.Add(1)
	go c.runClient()
}

// Stop stops the StreamClient.
func (c *StreamClient) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *StreamClient) runClient() {
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

func (c *StreamClient) runClientIteration() {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	conn, err := c.coordinatorClient.Connect(ctx)
	if err != nil {
		log.Autosample().Warn().Err(err).Msg("stream client: failed to connect")
		return
	}

	err = conn.Send(&rpcv1.ClientToServer{
		Msg: &rpcv1.ClientToServer_Hello_{
			Hello: &rpcv1.ClientToServer_Hello{
				Name: c.clientName,
				// We're creating a fresh dispatcher for every connection,
				// thus we can start from 0 each time.
				NextId: 0,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("stream client: failed to hello")
		return
	}

	dispatcher := c.handlers.StartDispatcher()
	defer dispatcher.Stop()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		for {
			// if ctx expires or conn disconnects, this should return error
			// and the loop should end, no need to explicitly wait on
			// ctx.done().
			serverToClient, err := conn.Recv()
			if err != nil {
				log.Warn().Err(err).Msg("stream client: disconnected")
				cancel()
				return
			}
			dispatcher.ProcessCommand(serverToClient)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case response := <-dispatcher.Chan():
			err := conn.Send(&rpcv1.ClientToServer{
				Msg: &rpcv1.ClientToServer_Response{
					Response: response,
				},
			})
			if err != nil {
				log.Warn().Err(err).Msg("stream client: failed to send")
				// Note: This cancel is deferred already, but we want to cancel
				// to force-stop serverToClient goroutine before joining it.
				cancel()
				return
			}
		}
	}
}

// StreamServer is a server part of "stream" transport, which uses
// bidirectional grpc stream methods.
//
// Implements rpc.v1.Coordinator.
type StreamServer struct {
	clients *Clients
}

// NewStreamServer creates new StreamServer plugging into given Clients as a transport.
func NewStreamServer(clients *Clients) *StreamServer {
	return &StreamServer{
		clients: clients,
	}
}

// Connect implements rpcv1.CoordinatorServer.
func (s *StreamServer) Connect(conn rpcv1.Coordinator_ConnectServer) error {
	firstMsg, err := conn.Recv()
	if err != nil {
		return err
	}

	hello, ok := firstMsg.Msg.(*rpcv1.ClientToServer_Hello_)
	if !ok {
		return status.Error(codes.InvalidArgument, "expected hello msg")
	}

	requests, responses := s.clients.Join(hello.Hello.Name, hello.Hello.NextId)
	defer close(responses)

	// Note: Capacity set to one, so goroutine would never block on sending to this channel.
	clientDisconnected := make(chan error, 1)

	go func() {
		for {
			msg, err := conn.Recv()
			if err != nil {
				if err == io.EOF {
					clientDisconnected <- nil
				} else {
					clientDisconnected <- err
				}
				return
			}

			switch msg := msg.Msg.(type) {
			case *rpcv1.ClientToServer_Response:
				responses <- msg.Response
			case *rpcv1.ClientToServer_Hello_:
				clientDisconnected <- status.Error(codes.InvalidArgument, "unexpected hello")
				return
			default:
				log.Bug().Msg("unknown client to server message")
			}
		}
	}()
	// Note: not waiting for this goroutine in defer, see below.

	for {
		select {
		case serverToClient, ok := <-requests:
			if !ok {
				// Close the stream.
				// Note: The goroutine should receive err from Recv and finish immediately.
				return nil
			}

			if err := conn.Send(serverToClient); err != nil {
				return err
			}

		case err := <-clientDisconnected:
			return err
		}
	}
}
