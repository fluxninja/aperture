package grpcui

import (
	"context"
	"net/http"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Module is a fx module that invokes gRPC UI handler and invokes registering gRPC UI handler.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(setupGRPCUIHandler),
	)
}

const (
	grpcUIEndpoint = "/grpcui"
)

func setupGRPCUIHandler(_ *grpc.Server, router *mux.Router, lifecycle fx.Lifecycle, unmarshaller config.Unmarshaller) error {
	var err error
	var conn *grpc.ClientConn
	var h http.Handler

	var serverConfig grpcclient.GRPCServerConfig
	err = unmarshaller.UnmarshalKey(grpcclient.DefaultServerConfigKey, &serverConfig)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deserialize grpc server configuration!")
		return err
	}

	var listenerConfig listener.ListenerConfig
	err = unmarshaller.UnmarshalKey("server", &listenerConfig)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deserialize listener configuration!")
		return err
	}

	targetAddr := listenerConfig.Addr

	// TODO: Remove Debug Log
	log.Debug().Interface("serverConfig", serverConfig).Msg("starting providegrpcuihandler -- server")
	log.Debug().Interface("listenerConfig", listenerConfig).Msg("starting providegrpcuihandler -- listener")

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if serverConfig.EnableGRPCUI {
				// TODO: Remove Debug Log
				log.Debug().Str("targetAddr", targetAddr).Msg("Providegrpcuuihandler target address")
				conn, err = grpc.DialContext(context.Background(), targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Error().Err(err).Msg("Failed to create gRPC ui client connection to the server")
					return err
				}
				// TODO: Connection closed before server preface received.
				h, err = standalone.HandlerViaReflection(context.Background(), conn, conn.Target())
				if err != nil {
					log.Error().Err(err).Msg("Failed to create gRPC ui handler")
					return err
				}
				router.Handle(grpcUIEndpoint+"/", http.StripPrefix(grpcUIEndpoint, h))
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})

	return nil
}
