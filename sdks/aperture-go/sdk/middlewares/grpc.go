package middlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	aperture "github.com/fluxninja/aperture-go/v2/sdk"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// socketAddressFromNetAddr takes a net.Addr and returns a flowcontrolhttp.SocketAddress.
func socketAddressFromNetAddr(logger logr.Logger, addr net.Addr) *flowcontrolhttp.SocketAddress {
	host, port := splitAddress(logger, addr.String())
	protocol := flowcontrolhttp.SocketAddress_TCP
	if addr.Network() == "udp" {
		protocol = flowcontrolhttp.SocketAddress_UDP
	}
	return &flowcontrolhttp.SocketAddress{
		Address:  host,
		Protocol: protocol,
		Port:     port,
	}
}

// GRPCUnaryInterceptor takes a control point name and creates a UnaryInterceptor which can be used with gRPC server.
func GRPCUnaryInterceptor(c aperture.Client, controlPoint string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		checkreq := PrepareCheckHTTPRequestForGRPC(req, ctx, info, c.GetLogger(), controlPoint)

		flow, err := c.StartHTTPFlow(ctx, checkreq)
		if err != nil {
			c.GetLogger().Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Accepted()", flow.Accepted())
		}

		defer func() {
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
			// The second argument is error message for further diagnosis.
			err := flow.End(aperture.OK)
			if err != nil {
				c.GetLogger().Info("Aperture flow control end got error.", "error", err)
			}
		}()

		if flow.Accepted() {
			// Simulate work being done
			resp, err := handler(ctx, req)
			return resp, err
		} else {
			rejectResp := flow.CheckResponse().GetDeniedResponse()
			return nil, status.Error(
				convertHTTPStatusToGRPC(rejectResp.GetStatus()),
				fmt.Sprintf("Aperture rejected the request: %v", rejectResp.GetBody()),
			)
		}
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
