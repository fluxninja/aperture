package rpc

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	rpcv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/rpc/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

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

// RegisterFunction register a function as ia handler in the registry
//
// Only one function for a given Req type can be registered.
//
// Note: This is not a method due to golang's generic's limitations.
func RegisterFunction[Req, Resp proto.Message](
	registry *HandlerRegistry,
	handler func(context.Context, Req) (Resp, error),
) error {
	var req Req // used only to pull out the message name from descriptor
	name := req.ProtoReflect().Descriptor().FullName()
	if _, exists := registry.handlers[name]; exists {
		return fmt.Errorf("handler for %s already registered", name)
	}

	registry.handlers[name] = func(ctx context.Context, anyreq *anypb.Any) (proto.Message, error) {
		var nilReq Req
		req := nilReq.ProtoReflect().New().Interface().(Req)
		if err := anyreq.UnmarshalTo(req); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return handler(ctx, req)
	}

	return nil
}

// StartDispatcher starts a dispatcher using registered handlers.
//
// Multiple independent dispatchers can be started using single handler registry.
func (r *HandlerRegistry) StartDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:       r.handlers,
		activeRequests: map[uint64]activeRequest{},
		responses:      make(chan *rpcv1.Response, 64),
	}
}

// Dispatcher dispatches ServerToClient calls and calls appropriate handler
//
// Dispatcher is intended to be used in a synchronized way from a *single*
// transport.
type Dispatcher struct {
	handlers            map[protoreflect.FullName]untypedHandler
	activeRequests      map[uint64]activeRequest
	activeRequestsMutex sync.Mutex
	responses           chan *rpcv1.Response
}

type activeRequest struct {
	cancel context.CancelFunc
}

// ProcessCommand processes a single server-to-client command.
func (d *Dispatcher) ProcessCommand(message *rpcv1.ServerToClient) {
	switch cmd := message.Command.(type) {

	case *rpcv1.ServerToClient_Request:
		ctx, cancel := context.WithTimeout(context.Background(), cmd.Request.Timeout.AsDuration())
		d.activeRequestsMutex.Lock()
		d.activeRequests[cmd.Request.Id] = activeRequest{cancel: cancel}
		d.activeRequestsMutex.Unlock()
		go d.callAndHandleResult(ctx, cmd.Request.Id, cmd.Request.Payload)

	case *rpcv1.ServerToClient_CancelId:
		d.activeRequestsMutex.Lock()
		req := d.activeRequests[cmd.CancelId]
		delete(d.activeRequests, cmd.CancelId)
		d.activeRequestsMutex.Unlock()
		if req.cancel != nil {
			req.cancel()
		}

	default:
		log.Bug().Msg("Unhandled ServerToClient command")
	}
}

// Chan returns a channel with responses that transport should forward to server.
func (d *Dispatcher) Chan() <-chan *rpcv1.Response { return d.responses }

// Stop cancels all the in-flight requests.
//
// Note: Doesn't wait for completion.
func (d *Dispatcher) Stop() {
	d.activeRequestsMutex.Lock()
	defer d.activeRequestsMutex.Unlock()

	for _, req := range d.activeRequests {
		req.cancel()
	}
	d.activeRequests = nil
}

func (d *Dispatcher) callAndHandleResult(ctx context.Context, id uint64, req *anypb.Any) {
	serializedResp, err := d.callHandler(ctx, req)

	if !d.deactivate(id) {
		// request already canceled, no need to send response
		return
	}

	if err != nil {
		d.responses <- &rpcv1.Response{
			Id: id,
			Result: &rpcv1.Response_Error{
				Error: status.Convert(err).Proto(),
			},
		}
	} else {
		d.responses <- &rpcv1.Response{
			Id: id,
			Result: &rpcv1.Response_Payload{
				Payload: serializedResp,
			},
		}
	}
}

func (d *Dispatcher) callHandler(ctx context.Context, req *anypb.Any) ([]byte, error) {
	handler, exists := d.handlers[req.MessageName()]
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

func (d *Dispatcher) deactivate(id uint64) bool {
	d.activeRequestsMutex.Lock()
	defer d.activeRequestsMutex.Unlock()

	_, wasActive := d.activeRequests[id]
	if wasActive {
		delete(d.activeRequests, id)
	}
	return wasActive
}
