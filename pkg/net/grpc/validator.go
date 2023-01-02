package grpc

import (
	"context"

	"github.com/fluxninja/aperture/pkg/config"
	"google.golang.org/grpc"
)

// validatorUnaryInterceptor returns a new unary server interceptors that applies defaults and validates the request
// before calling the handler.
func validatorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// set defaults for the req
		config.SetDefaults(req)
		// validate the req
		if err := config.ValidateStruct(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
