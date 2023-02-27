// Server-to-client reverse RPC
//
// Idea: clients connect to server and server can call functions on clients
// (either on specified client or on all simultaneously).
// This is mainly for purpose of controller to call functions on agents, but
// can be used for different purposes too.
//
// The implementation is modularized as follows:
//
//	server      [transport] <-> [Clients]    <- callers
//	                 ^
//	-----------------|----------------------------------
//	                 v
//	client      [transport] <-> [Dispatcher] -> handlers (functions)
//
// Clients manages connected clients and provides an api for listing clients
// and calling functions.
//
// Dispatcher dispatches messages from transport to specific handlers.
//
// Handlers are simply functions `func(ctx, req) (resp, error)` that can be
// registered in HandlerRegistry.
//
// Transports (client & server part) provide "actual transport IO".
//
// This is a bit m:1:n architecture:
// * Dispatcher is a bridge between transports and handlers.
// * Clients is a bridge between transport and callers.
// Thus implementing additional handlers or transports should not require any
// changes in Client nor Dispatcher.
//
// Transports should only be concerned in connecting and passing ClientToServer
// and ServerToClient messages. Currently implemented transports are:
//   - Bidirectional-grpc-stream-based (rpc.v1.Coordinator.Connect).
//
// As grpc bidirectional stream may not work everywhere, possible and planned
// transports are:
//   - Heartbeat-based (heartbeat request passes ClientToServer messages,
//     response contains ServerToClient messages; some additional latency is
//     introduced).
//   - Long-poll.
package rpc
