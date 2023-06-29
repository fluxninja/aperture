package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	aperture "github.com/fluxninja/aperture-go/v2/sdk"
)

// splits address into host and port.
func splitAddress(logger logr.Logger, address string) (string, uint32) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		logger.V(2).Info("Error splitting address", "address", address, "error", err)
		return host, 0
	}

	portUint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		logger.V(2).Info("Error parsing port", "address", address, "error", err)
		return host, 0
	}

	return host, uint32(portUint)
}

// reads body from request, replacing it with a clone to allow further reads.
func readClonedBody(r *http.Request) ([]byte, error) {
	body := r.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// try to figure out the local ip address.
func getLocalIP(logger logr.Logger) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.V(2).Info("Failed to get local IP address", "error", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	logger.V(2).Info("Failed to get local IP address")
	return ""
}

// PrepareCheckHTTPRequestForHTTP takes a http.Request, logger and Control Point to use in Aperture policy for preparing the flowcontrolhttp.CheckHTTPRequest and returns it.
func PrepareCheckHTTPRequestForHTTP(httpRequest *http.Request, logger logr.Logger, controlPoint string) *flowcontrolhttp.CheckHTTPRequest {
	labels := aperture.LabelsFromCtx(httpRequest.Context())

	for key, value := range httpRequest.Header {
		if strings.HasPrefix(key, ":") {
			continue
		}
		labels[key] = strings.Join(value, ",")
	}

	// We know that the protocol is TCP because Golang's http package doesn't support UDP
	// TODO: Should we support `httpu`?
	protocol := flowcontrolhttp.SocketAddress_TCP

	sourceHost, sourcePort := splitAddress(logger, httpRequest.RemoteAddr)
	// TODO: Figure out if we can narrow down the port or figure out the host in a better way
	destinationPort := uint32(0)
	destinationHost := getLocalIP(logger)

	bodyBytes, err := readClonedBody(httpRequest)
	if err != nil {
		logger.V(2).Info("Error reading body", "error", err)
	}

	return &flowcontrolhttp.CheckHTTPRequest{
		Source: &flowcontrolhttp.SocketAddress{
			Address:  sourceHost,
			Protocol: protocol,
			Port:     sourcePort,
		},
		Destination: &flowcontrolhttp.SocketAddress{
			Address:  destinationHost,
			Protocol: protocol,
			Port:     destinationPort,
		},
		ControlPoint: controlPoint,
		Request: &flowcontrolhttp.CheckHTTPRequest_HttpRequest{
			Method:   httpRequest.Method,
			Path:     httpRequest.URL.Path,
			Host:     httpRequest.Host,
			Headers:  labels,
			Scheme:   httpRequest.URL.Scheme,
			Size:     httpRequest.ContentLength,
			Protocol: httpRequest.Proto,
			Body:     string(bodyBytes),
		},
	}
}

// PrepareCheckHTTPRequestForGRPC takes a gRPC request, context, unary server-info, logger and Control Point to use in Aperture policy for preparing the flowcontrolhttp.CheckHTTPRequest and returns it.
func PrepareCheckHTTPRequestForGRPC(req interface{}, ctx context.Context, info *grpc.UnaryServerInfo, logger logr.Logger, controlPoint string) *flowcontrolhttp.CheckHTTPRequest {
	labels := aperture.LabelsFromCtx(ctx)

	md, ok := metadata.FromIncomingContext(ctx)
	authority := ""
	scheme := ""
	method := ""

	if ok {
		for key, value := range md {
			labels[key] = strings.Join(value, ",")
		}
		getMetaValue := func(key string) string {
			values := md.Get(key)
			if len(values) > 0 {
				return values[0]
			}
			return ""
		}
		authority = getMetaValue(":authority")
		scheme = getMetaValue(":scheme")
		method = getMetaValue(":method")
	}

	var sourceSocket *flowcontrolhttp.SocketAddress
	if sourceAddr, ok := peer.FromContext(ctx); ok {
		sourceSocket = socketAddressFromNetAddr(logger, sourceAddr.Addr)
	}
	destinationSocket := &flowcontrolhttp.SocketAddress{
		Address:  getLocalIP(logger),
		Protocol: flowcontrolhttp.SocketAddress_TCP,
		Port:     0,
	}

	body, err := json.Marshal(req)
	if err != nil {
		logger.V(2).Info("Unable to marshal request body")
	}

	return &flowcontrolhttp.CheckHTTPRequest{
		Source:       sourceSocket,
		Destination:  destinationSocket,
		ControlPoint: controlPoint,
		Request: &flowcontrolhttp.CheckHTTPRequest_HttpRequest{
			Method:   method,
			Path:     info.FullMethod,
			Host:     authority,
			Headers:  labels,
			Scheme:   scheme,
			Size:     -1,
			Protocol: "HTTP/2",
			Body:     string(body),
		},
	}
}
