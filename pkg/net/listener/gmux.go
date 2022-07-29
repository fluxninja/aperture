package listener

import (
	"net/http"

	"github.com/elastic/gmux"
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/log"
)

// GMuxConstructor holds fields to create an annotated instance gmux Listener.
type GMuxConstructor struct {
	// Name of HTTP Server instance
	HTTPServerName string
	// Name of GRPC listener instance
	ListenerName string
}

// Annotate creates an annotated instance of gmux Listener.
func (constructor GMuxConstructor) Annotate() fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor.configureServer,
			fx.ParamTags(config.NameTag(constructor.HTTPServerName)),
			fx.ResultTags(config.NameTag(constructor.ListenerName)),
		),
	)
}

func (constructor GMuxConstructor) configureServer(server *http.Server) (*Listener, error) {
	listener := &Listener{}
	var err error
	// Listener is automatically closed by http server on shutdown
	listener.lis, err = gmux.ConfigureServer(server, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to setup gmux listener!")
		return nil, err
	}

	return listener, nil
}
