package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// socketAddressFromNetAddr takes a net.Addr and returns a flowcontrolhttp.SocketAddress.
func (c *apertureClient) socketAddressFromNetAddr(addr net.Addr) *flowcontrolhttp.SocketAddress {
	host, port := c.splitAddress(addr.String())
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
func (c *apertureClient) GRPCUnaryInterceptor(controlPoint string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		labels := labelsFromCtx(ctx)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for key, value := range md {
				labels[key] = strings.Join(value, ",")
			}
		}

		var sourceSocket *flowcontrolhttp.SocketAddress
		if sourceAddr, ok := peer.FromContext(ctx); ok {
			sourceSocket = c.socketAddressFromNetAddr(sourceAddr.Addr)
		}
		// TODO: Can we retrieve the address somehow?
		var destinationSocket *flowcontrolhttp.SocketAddress
		//if server, ok := info.Server.(*grpc.Server); ok {
		//	service := strings.Split(path.Dir(info.FullMethod), "/")[1]
		//	if serviceInfo, ok := server.GetServiceInfo()[service]; ok {
		//		if listener, ok := serviceInfo.Metadata["listener"]; ok {
		//			if listener2, ok := listener.(net.Listener); ok {
		//				destinationSocket = c.socketAddressFromNetAddr(listener2.Addr())
		//			}
		//		}
		//	}
		//}

		// TODO: Fill this up properly
		checkreq := &flowcontrolhttp.CheckHTTPRequest{
			Source:       sourceSocket,
			Destination:  destinationSocket,
			ControlPoint: controlPoint,
			Request: &flowcontrolhttp.CheckHTTPRequest_HttpRequest{
				Method:   "",
				Path:     info.FullMethod,
				Host:     "",
				Headers:  labels,
				Scheme:   "",
				Size:     -1,
				Protocol: "HTTP/2",
				Body:     json.Marshal(req),
			},
		}

		flow, err := c.StartHTTPFlow(ctx, checkreq)
		if err != nil {
			c.log.Info("Aperture flow control got error. Returned flow defaults to Allowed.", "flow.Accepted()", flow.Accepted())
		}

		if flow.Accepted() {
			// Simulate work being done
			resp, err := handler(ctx, req)
			// Need to call End() on the Flow in order to provide telemetry to Aperture Agent for completing the control loop.
			// The first argument captures whether the feature captured by the Flow was successful or resulted in an error.
			// The second argument is error message for further diagnosis.
			flowErr := flow.End(OK)
			if flowErr != nil {
				c.log.Info("Aperture flow control end got error.", "error", err)
			}
			return resp, err
		} else {
			err := flow.End(OK)
			if err != nil {
				c.log.Info("Aperture flow control end got error.", "error", err)
			}
			rejectResp := flow.CheckResponse().GetRejectedResponse()
			return nil, status.Error(
				rejectResp.GetStatus(),
				fmt.Sprintf("Aperture rejected the request: %v", rejectResp.GetBody()),
			)
		}
	}
}
