package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// GRPCUnaryInterceptor takes a control point name and creates a UnaryInterceptor which can be used with gRPC server.
func GRPCUnaryInterceptor(c aperture.Client, controlPoint string, explicitLabels map[string]string, rampMode bool, timeout time.Duration) grpc.UnaryServerInterceptor {
	if explicitLabels == nil {
		explicitLabels = make(map[string]string)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		checkreq, err := prepareCheckHTTPRequestForGRPC(req, ctx, info, controlPoint, explicitLabels, rampMode)
		if err != nil {
			c.GetLogger().Error("Failed to prepare CheckHTTP request.", "error", err)
		}

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

// PrepareCheckHTTPRequestForGRPC takes a gRPC request, context, unary server-info, logger and Control Point to use in Aperture policy for preparing the flowcontrolhttp.CheckHTTPRequest and returns it.
func prepareCheckHTTPRequestForGRPC(req interface{}, ctx context.Context, info *grpc.UnaryServerInfo, controlPoint string, explicitLabels map[string]string, rampMode bool) (*checkhttpv1.CheckHTTPRequest, error) {
	labels := utils.LabelsFromCtx(ctx)

	// override labels with explicit labels
	for key, value := range explicitLabels {
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
		return nil, err
	}

	return &checkhttpv1.CheckHTTPRequest{
		Source:       sourceSocket,
		Destination:  destinationSocket,
		ControlPoint: controlPoint,
		RampMode:     rampMode,
		Request: &checkhttpv1.CheckHTTPRequest_HttpRequest{
			Method:   method,
			Path:     info.FullMethod,
			Host:     authority,
			Headers:  labels,
			Scheme:   scheme,
			Size:     -1,
			Protocol: "HTTP/2",
			Body:     string(body),
		},
	}, nil
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
