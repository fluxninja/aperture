package grpcui

import (
	"context"
	"net/http"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

// Module is a fx module that provides gRPC UI handler and invokes registering gRPC UI handler.
func Module() fx.Option {
	return fx.Options(
		grpcclient.ClientConstructor{Name: "grpcui-grpc-client", ConfigKey: "grpcui-grpc-client"}.Annotate(),
		fx.Provide(fx.Annotate(
			provideGRPCUIHandler,
			fx.ParamTags(config.NameTag("grpcui-grpc-client")),
		)),
		fx.Invoke(RegisterGRPCUIHandler),
	)
}

const (
	grpcUIEndpoint = "/grpcui"
)

func provideGRPCUIHandler(GRPClientConnectionBuilder grpcclient.ClientConnectionBuilder, HTTPServer *commonhttp.Server, Unmarshaller config.Unmarshaller, Lifecycle fx.Lifecycle) (http.Handler, error) {
	var serverConfig grpcclient.GRPCServerConfig
	err := Unmarshaller.UnmarshalKey(grpcclient.DefaultServerConfigKey, &serverConfig)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deserialize grpc server configuration!")
		return nil, err
	}

	var h http.Handler
	Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if serverConfig.EnableGRPCUI {
				connWrapper := GRPClientConnectionBuilder.Build()
				targetAddr := HTTPServer.Listener.GetListener().Addr().String()
				// TODO: Do not disable transport security.
				conn, err := connWrapper.Dial(ctx, targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Error().Err(err).Msg("Failed to create gRPC ui client connection to the server")
					return err
				}
				// TODO: Connection closed before server preface received.
				h, err = standalone.HandlerViaReflection(ctx, conn, conn.Target())
				if err != nil {
					log.Error().Err(err).Msg("Failed to create gRPC ui handler")
					return err
				}
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return h, nil
}

// RegisterGRPCUIHandler registers gRPC UI handler to router.
func RegisterGRPCUIHandler(router *mux.Router, handler http.Handler) {
	log.Info().Msg("Registering gRPC UI handler")
	router.Handle(grpcUIEndpoint, http.StripPrefix(grpcUIEndpoint, handler))
}
