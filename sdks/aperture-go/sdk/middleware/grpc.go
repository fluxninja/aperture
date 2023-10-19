package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

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

// GRPCUnaryInterceptor takes a control point name and creates a UnaryInterceptor which can be used with gRPC server.
func GRPCUnaryInterceptor(c aperture.Client, controlPoint string, explicitLabels map[string]string, rampMode bool, timeout time.Duration) grpc.UnaryServerInterceptor {
	if explicitLabels == nil {
		explicitLabels = make(map[string]string)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		checkreq := prepareCheckHTTPRequestForGRPC(req, ctx, info, c.GetLogger(), controlPoint, explicitLabels, rampMode)

		flow := c.StartHTTPFlow(ctx, checkreq, rampMode, timeout)
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

		if flow.ShouldRun() {
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
