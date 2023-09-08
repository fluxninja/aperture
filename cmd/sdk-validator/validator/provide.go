package validator

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Module .
var Module = fx.Options(
	fx.Provide(ProvideFlowControlHandler),
	fx.Invoke(Register),
)

// ConstructorIn .
type ConstructorIn struct {
	fx.In
}

// ProvideFlowControlHandler .
func ProvideFlowControlHandler(in ConstructorIn) (flowcontrolv1.FlowControlServiceServer, error) {
	return &FlowControlHandler{}, nil
}

// RegisterIn bundles and annotates parameters.
type RegisterIn struct {
	fx.In
	Server  *grpc.Server `name:"default"`
	Handler flowcontrolv1.FlowControlServiceServer
}

// Register registers flowcontrol service on a gRPC server.
func Register(in RegisterIn) {
	flowcontrolv1.RegisterFlowControlServiceServer(in.Server, in.Handler)
	log.Info().Msg("flowcontrol handler registered")
}
