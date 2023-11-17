package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	checkhttpv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	aperture "github.com/fluxninja/aperture-go/v2/sdk"
	"github.com/fluxninja/aperture-go/v2/sdk/utils"
)

// socketAddressFromNetAddr takes a net.Addr and returns a flowcontrolhttp.SocketAddress.
func socketAddressFromNetAddr(addr net.Addr) *checkhttpv1.SocketAddress {
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil
	}

	portU32, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return nil
	}

	protocol := checkhttpv1.SocketAddress_TCP
	if addr.Network() == "udp" {
		protocol = checkhttpv1.SocketAddress_UDP
	}
	return &checkhttpv1.SocketAddress{
		Address:  host,
		Protocol: protocol,
		Port:     uint32(portU32),
	}
}

// NewGRPCMiddleware takes a control point name and creates a UnaryInterceptor which can be used with gRPC server.
func NewGRPCMiddleware(client aperture.Client, controlPoint string, middlewareParams aperture.MiddlewareParams) (grpc.UnaryServerInterceptor, error) {
	// Precompile the regex patterns for ignored paths
	if middlewareParams.IgnoredPaths != nil {
		compiledIgnoredPaths := make([]*regexp.Regexp, len(middlewareParams.IgnoredPaths))
		for i, pattern := range middlewareParams.IgnoredPaths {
			compiledPattern, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			} else {
				compiledIgnoredPaths[i] = compiledPattern
			}
		}
		middlewareParams.IgnoredPathsCompiled = compiledIgnoredPaths
	}

	return GRPCUnaryInterceptor(client, controlPoint, middlewareParams), nil
}

// GRPCUnaryInterceptor takes a control point name and creates a UnaryInterceptor which can be used with gRPC server.
func GRPCUnaryInterceptor(c aperture.Client, controlPoint string, middlewareParams aperture.MiddlewareParams) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// If the path is ignored, skip the middleware
		if middlewareParams.IgnoredPathsCompiled != nil {
			for _, ignoredPath := range middlewareParams.IgnoredPathsCompiled {
				if ignoredPath.MatchString(info.FullMethod) {
					return handler(ctx, req)
				}
			}
		}

		checkReq := prepareCheckHTTPRequestForGRPC(ctx, req, c.GetLogger(), info.FullMethod, controlPoint, middlewareParams.FlowParams)

		flow := c.StartHTTPFlow(ctx, checkReq, middlewareParams)
		if flow.Error() != nil {
			c.GetLogger().Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Error()", flow.Error().Error(), "flow.ShouldRun()", flow.ShouldRun())
		}

		defer func() {
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// SetStatus() method of Flow object can be used to capture whether the Flow was successful or resulted in an error.
			// If not set, status defaults to OK.
			err := flow.End()
			if err != nil {
				c.GetLogger().Info("Aperture flow control end got error.", "error", err)
			}
		}()

		if !flow.ShouldRun() {
			rejectResp := flow.CheckResponse().GetDeniedResponse()
			return nil, status.Error(
				convertHTTPStatusToGRPC(rejectResp.GetStatus()),
				fmt.Sprintf("Aperture rejected the request: %v", rejectResp.GetBody()),
			)
		}

		return handler(ctx, req)
	}
}

// PrepareCheckHTTPRequestForGRPC takes a gRPC request, context, unary server-info, logger and Control Point to use in Aperture policy for preparing the flowcontrolhttp.CheckHTTPRequest and returns it.
func prepareCheckHTTPRequestForGRPC(ctx context.Context, req interface{}, logger *slog.Logger, fullMethod string, controlPoint string, flowParams aperture.FlowParams) *checkhttpv1.CheckHTTPRequest {
	labels := utils.LabelsFromCtx(ctx)

	// override labels with explicit labels
	for key, value := range flowParams.Labels {
		labels[key] = value
	}

	md, ok := metadata.FromIncomingContext(ctx)
	authority := ""
	scheme := ""
	method := ""

	if ok {
		// override labels with labels from metadata
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

	var sourceSocket *checkhttpv1.SocketAddress
	if sourceAddr, ok := peer.FromContext(ctx); ok {
		sourceSocket = socketAddressFromNetAddr(sourceAddr.Addr)
	}
	destinationSocket := &checkhttpv1.SocketAddress{
		Address:  utils.GetLocalIP(),
		Protocol: checkhttpv1.SocketAddress_TCP,
		Port:     0,
	}

	body, err := json.Marshal(req)
	if err != nil {
		logger.Error("Failed to marshal request body", "error", err)
	}

	return &checkhttpv1.CheckHTTPRequest{
		Source:       sourceSocket,
		Destination:  destinationSocket,
		ControlPoint: controlPoint,
		RampMode:     flowParams.RampMode,
		Request: &checkhttpv1.CheckHTTPRequest_HttpRequest{
			Method:   method,
			Path:     fullMethod,
			Host:     authority,
			Headers:  labels,
			Scheme:   scheme,
			Size:     -1,
			Protocol: "HTTP/2",
			Body:     string(body),
		},
	}
}

func convertHTTPStatusToGRPC(httpStatusCode int32) codes.Code {
	switch httpStatusCode {
	case http.StatusOK:
		return codes.OK
	case http.StatusRequestTimeout:
		return codes.Canceled
	case http.StatusInternalServerError:
		return codes.Unknown
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusPreconditionFailed:
		return codes.FailedPrecondition
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}
