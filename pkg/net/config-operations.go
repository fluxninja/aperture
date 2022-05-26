package net

// Dummy Swagger operations that represent yaml based configuration. For use with markdown generator only.

// swagger:operation POST /server common-configuration Server
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ListenerConfig"
// - name: http
//   in: body
//   schema:
//     "$ref": "#/definitions/HTTPServerConfig"
// - name: grpc
//   in: body
//   schema:
//     "$ref": "#/definitions/GRPCServerConfig"
// - name: grpc_gateway
//   in: body
//   schema:
//     "$ref": "#/definitions/GRPCGatewayConfig"
// - name: tls
//   in: body
//   schema:
//     "$ref": "#/definitions/ServerTLSConfig"

// swagger:operation POST /client common-configuration Client
// ---
// x-fn-config-env: true
// parameters:
// - name: proxy
//   in: body
//   schema:
//     "$ref": "#/definitions/ProxyConfig"
