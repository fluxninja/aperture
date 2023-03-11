package rpc

import (
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	rpcv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/rpc/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

const defaultTimeout = 5 * time.Second

// Clients manages connected clients and allows calling functions on them.
type Clients struct {
	clients         map[string]*connectedClient
	nextGlobalReqID uint64
	mutex           sync.Mutex
}

type connectedClient struct {
	requests chan request
	// Whether the client is already isDeleted from the map.
	// Protected by Clients.Mutex.
	isDeleted bool
	// Invariant: request chan is closed iff isDeleted
}

// A request or its cancellation.
type request struct {
	globalReqID uint64           // request ID from caller's POV
	req         *anypb.Any       // nil to cancel in-flight request
	respChan    chan<- RawResult // applicable when req != nil
	timeout     time.Duration    // applicable when req != nil
	// applicable if req == nil. If the request is canceled because of timeout,
	// no need to send cancellation request, as it should time out on its own.
	// If cancellation is "manual", cancellation request should be sent.
	isTimeout bool
}

// NewClients creates new Clients.
func NewClients() *Clients {
	return &Clients{
		clients: map[string]*connectedClient{},
	}
}

// APIs for transports -------------

// Join marks a new client as connected.
//
// When client disconnects, transport should close responses channel.
//
// ServerToClient channel will remain unclosed unless client with the same name will Join.
// In such case, transport should disconnect the client, as such connection is "stale".
// (Note: in this case transport should close responses channel as usual too).
func (c *Clients) Join(
	clientName string,
	nextID uint64,
) (<-chan *rpcv1.ServerToClient, chan<- *rpcv1.Response) {
	commands := make(chan *rpcv1.ServerToClient, 1)
	responses := make(chan *rpcv1.Response, 1)
	requests := make(chan request, 1)

	log.Info().Str("name", clientName).Msg("rpc: Client joined")

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if prevClient, exists := c.clients[clientName]; exists {
		close(prevClient.requests)
		prevClient.isDeleted = true
		delete(c.clients, clientName)
	}

	client := &connectedClient{requests: requests}
	c.clients[clientName] = client
	go c.handleClient(clientName, nextID, client, commands, responses)

	return commands, responses
}

func (c *Clients) handleClient(
	name string,
	nextID uint64,
	client *connectedClient,
	commands chan<- *rpcv1.ServerToClient,
	responses <-chan *rpcv1.Response,
) {
	responseChans := map[uint64]chan<- RawResult{}
	localToGlobalID := map[uint64]uint64{}
	globalToLocalID := map[uint64]uint64{}

	free := func(id, globalID uint64) {
		delete(responseChans, id)
		delete(localToGlobalID, id)
		delete(globalToLocalID, globalID)
	}

	handleClientResponse := func(resp *rpcv1.Response) {
		respChan, exists := responseChans[resp.Id]
		if !exists {
			return
		}
		free(resp.Id, localToGlobalID[resp.Id])

		switch result := resp.Result.(type) {
		case *rpcv1.Response_Payload:
			respChan <- RawResult{
				Client:  name,
				Success: result.Payload,
			}
		case *rpcv1.Response_Error:
			respChan <- RawResult{
				Client: name,
				Err:    status.ErrorProto(result.Error),
			}
		}
	}

	handleClientDisconnect := func() {
		for _, respChan := range responseChans {
			respChan <- RawResult{
				Client: name,
				Err:    status.Error(codes.Canceled, "client disconnected"),
			}
		}

		c.mutex.Lock()
		defer c.mutex.Unlock()
		if !client.isDeleted {
			delete(c.clients, name)
		}
	}

	for {
		select {
		case req, ok := <-client.requests:
			if !ok {
				// client kicked out by Clients
				close(commands)
				// Wait for client to disconnect by draining its channel.
				for resp := range responses {
					handleClientResponse(resp)
				}
				handleClientDisconnect()
				return
			}

			if req.req == nil {
				// Cancellation request
				id := globalToLocalID[req.globalReqID]
				if !req.isTimeout {
					commands <- &rpcv1.ServerToClient{
						Command: &rpcv1.ServerToClient_CancelId{
							CancelId: id,
						},
					}
				}
				free(id, req.globalReqID)
				break
			}

			id := nextID
			nextID += 1
			responseChans[id] = req.respChan
			localToGlobalID[id] = req.globalReqID
			globalToLocalID[req.globalReqID] = id

			commands <- &rpcv1.ServerToClient{
				Command: &rpcv1.ServerToClient_Request{
					Request: &rpcv1.Request{
						Id:      id,
						Timeout: durationpb.New(req.timeout),
						Payload: req.req,
					},
				},
			}

		case resp, ok := <-responses:
			if !ok {
				handleClientDisconnect()
				return
			}
			handleClientResponse(resp)
		}
	}
}

// APIs for client "users" ---------

// List lists all active clients.
func (c *Clients) List() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	clients := make([]string, 0, len(c.clients))
	for name := range c.clients {
		clients = append(clients, name)
	}
	return clients
}

// Call calls given client.
func (c *Clients) Call(clientName string, req *anypb.Any) ([]byte, error) {
	timeout := defaultTimeout
	respChan := make(chan RawResult, 1)

	c.mutex.Lock()
	id := c.nextGlobalReqID
	c.nextGlobalReqID += 1
	// Note: sending request is done inside the lock to avoid sending on closed channel.
	client, exists := c.clients[clientName]
	if exists {
		client.requests <- request{
			globalReqID: id,
			req:         req,
			respChan:    respChan,
			timeout:     timeout,
		}
	}
	c.mutex.Unlock()

	if !exists {
		return nil, status.Error(codes.NotFound, "client not found")
	}

	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		c.cancelRequest(client, id, true)
		return nil, status.Error(codes.DeadlineExceeded, "timeout")
	case result := <-respChan:
		timer.Stop()
		return result.Success, nil
	}
}

// Result is a result from one client. Used in CallAll.
type Result[Resp any] struct {
	Client  string
	Success Resp
	Err     error
}

// RawResult is a Result that's not yet converted to the proper response type.
type RawResult = Result[[]byte]

// CallAll calls all clients and returns responses from all of them (in arbitrary order).
func (c *Clients) CallAll(req *anypb.Any) []RawResult {
	timeout := defaultTimeout

	c.mutex.Lock()
	respChan := make(chan RawResult, len(c.clients))
	remainingClients := make(map[string]*connectedClient, cap(respChan))
	id := c.nextGlobalReqID
	c.nextGlobalReqID += 1
	for clientName, client := range c.clients {
		client.requests <- request{
			globalReqID: id,
			req:         req,
			respChan:    respChan,
			timeout:     timeout,
		}
		remainingClients[clientName] = client
	}
	c.mutex.Unlock()

	if len(remainingClients) == 0 {
		return nil
	}

	timer := time.NewTimer(timeout)
	results := make([]RawResult, 0, len(remainingClients))
	for {
		select {
		case <-timer.C:
			for clientName, client := range remainingClients {
				c.cancelRequest(client, id, true)
				results = append(results, RawResult{
					Client: clientName,
					Err:    status.Error(codes.DeadlineExceeded, "timeout"),
				})
			}
			return results
		case result := <-respChan:
			results = append(results, result)
			delete(remainingClients, result.Client)
			if len(remainingClients) == 0 {
				timer.Stop()
				return results
			}
		}
	}
}

func (c *Clients) cancelRequest(client *connectedClient, globalReqID uint64, isTimeout bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !client.isDeleted {
		client.requests <- request{
			globalReqID: globalReqID,
			isTimeout:   isTimeout,
		}
	}
}

// Call is Clients.Call with conversion to typed response
//
// Intended as a building-block for typed wrappers.
func Call[RespValue any, Resp ExactMessage[RespValue]](clients *Clients, client string, req proto.Message) (*RespValue, error) {
	anyreq, err := anypb.New(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	respBytes, err := clients.Call(client, anyreq)
	if err != nil {
		return nil, err
	}

	var resp RespValue
	if err := proto.Unmarshal(respBytes, Resp(&resp)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &resp, nil
}

// CallAll is Clients.CallAll with conversion to typed response
//
// Intended as a building-block for typed wrappers.
func CallAll[RespValue any, Resp ExactMessage[RespValue]](clients *Clients, req proto.Message) ([]Result[*RespValue], error) {
	anyreq, err := anypb.New(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rawResults := clients.CallAll(anyreq)

	resps := make([]Result[*RespValue], 0, len(rawResults))
	for _, rawResult := range rawResults {
		if rawResult.Err != nil {
			resps = append(resps, Result[*RespValue]{
				Client: rawResult.Client,
				Err:    rawResult.Err,
			})
		} else {
			result := Result[*RespValue]{
				Client:  rawResult.Client,
				Success: new(RespValue),
			}
			if err := proto.Unmarshal(rawResult.Success, Resp(result.Success)); err != nil {
				result.Success = nil
				result.Err = status.Error(codes.InvalidArgument, err.Error())
			}
			resps = append(resps, result)
		}
	}

	return resps, nil
}

// ExactMessage is like proto.Message, but known to point to T
//
// See
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example
type ExactMessage[T any] interface {
	proto.Message
	*T
}
