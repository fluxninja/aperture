package health

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
)

// HealthModule is a module that provides grpc health server for checking services status.
func HealthModule() fx.Option {
	return fx.Options(
		grpcclient.ClientConstructor{Name: "health-grpc-client", Key: "health.client.grpc"}.Annotate(),
		fx.Provide(ProvideHealthServer),
		fx.Provide(fx.Annotate(
			ProvideHealthClient,
			fx.ParamTags(config.NameTag("health-grpc-client")),
		)),
		fx.Invoke(RegisterHealthServer),
	)
}

// ProvideHealthServer creates instance of health server.
func ProvideHealthServer(lifecycle fx.Lifecycle) *health.Server {
	server := health.NewServer()

	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			server.Shutdown()
			return nil
		},
	})

	return server
}

// ProvideHealthClient creates instance of client to health server.
func ProvideHealthClient(GRPClientConnectionBuilder grpcclient.ClientConnectionBuilder) (grpc_health_v1.HealthClient, error) {
	// Setup connection to health service
	connWrapper := GRPClientConnectionBuilder.Build()
	conn, err := connWrapper.Dial(context.Background(), "localhost:80")
	if err != nil {
		log.Warn().Err(err).Msg("Could not connect to health grpc server")
		return nil, err
	}

	healthClient := grpc_health_v1.NewHealthClient(conn)
	return healthClient, nil
}

// RegisterHealthServer registers health server to grpc_health_v1 api and sets default statuses.
func RegisterHealthServer(srv *grpc.Server, healthsrv *health.Server) {
	// It registers empty name server implicitly
	grpc_health_v1.RegisterHealthServer(srv, healthsrv)
}
