package validator

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// Module .
var Module = fx.Options(
	fx.Provide(ProvideFlowcontrolHandler),
	fx.Invoke(Register),
)

// ConstructorIn .
type ConstructorIn struct {
	fx.In
}

// ProvideFlowcontrolHandler .
func ProvideFlowcontrolHandler(in ConstructorIn) (flowcontrolv1.FlowControlServiceServer, error) {
	return &FlowControlHandler{}, nil
}

// Register registers flowcontrol service on a grpc server.
func Register(server *grpc.Server, handler flowcontrolv1.FlowControlServiceServer) {
	flowcontrolv1.RegisterFlowControlServiceServer(server, handler)
	log.Info().Msg("flowcontrol handler registered")
}
