package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"

	checkhttpproto "buf.build/gen/go/fluxninja/aperture/protocolbuffers/go/aperture/flowcontrol/checkhttp/v1"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	aperture "github.com/fluxninja/aperture-go/v2/sdk"
)

// socketAddressFromNetAddr takes a net.Addr and returns a flowcontrolhttp.SocketAddress.
func socketAddressFromNetAddr(logger logr.Logger, addr net.Addr) *checkhttpproto.SocketAddress {
	host, port := splitAddress(logger, addr.String())
	protocol := checkhttpproto.SocketAddress_TCP
	if addr.Network() == "udp" {
		protocol = checkhttpproto.SocketAddress_UDP
	}
	return &checkhttpproto.SocketAddress{
		Address:  host,
		Protocol: protocol,
		Port:     port,
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

		checkReq := prepareCheckHTTPRequestForGRPC(req, ctx, info, c.GetLogger(), controlPoint, middlewareParams.FlowParams)

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
