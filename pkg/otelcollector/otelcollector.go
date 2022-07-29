package otelcollector

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/converter/overwritepropertiesconverter"
	"go.opentelemetry.io/collector/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panic"
)

const schemeName = "file"

// Module is a fx module that invokes OTEL Collector.
func Module() fx.Option {
	return fx.Invoke(Invoke)
}

// ConstructorIn describes parameters passed to create OTEL Collector, server providing the OpenTelemetry Collector service.
type ConstructorIn struct {
	fx.In
	Lifecycle     fx.Lifecycle
	Shutdowner    fx.Shutdowner
	Factories     component.Factories
	Unmarshaller  config.Unmarshaller
	BaseConfig    *OTELConfig   `name:"base"`
	PluginConfigs []*OTELConfig `group:"plugin"`
}

// Invoke creates and runs a new instance of OTEL Collector with the passed configuration.
func Invoke(in ConstructorIn) (*service.Collector, error) {
	locations := []string{"file:main"}
	mapProviders := map[string]confmap.Provider{
		"file": NewOTELConfigUnmarshaler(in.BaseConfig.AsMap()),
	}
	for i, pluginConfig := range in.PluginConfigs {
		scheme := fmt.Sprintf("plugin-%v", i)
		locations = append(locations, fmt.Sprintf("%v:%v", scheme, scheme))
		mapProviders[scheme] = NewOTELConfigUnmarshaler(pluginConfig.AsMap())
	}

	configProvider, err := service.NewConfigProvider(service.ConfigProviderSettings{
		Locations:    locations,
		MapProviders: mapProviders,
		MapConverters: []confmap.Converter{
			overwritepropertiesconverter.New([]string{}),
			expandconverter.New(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating OTEL config provider: %w", err)
	}
	otelService, err := service.New(
		service.CollectorSettings{
			BuildInfo:               component.NewDefaultBuildInfo(),
			Factories:               in.Factories,
			ConfigProvider:          configProvider,
			DisableGracefulShutdown: true,
			LoggingOptions: []zap.Option{zap.WrapCore(func(zapcore.Core) zapcore.Core {
				return log.NewZerologAdapter(log.GetGlobalLogger())
			})},
			// NOTE: do not remove this becauase it causes a data-race condition.
			SkipSettingGRPCLogger: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("constructing OTEL Service: %v", err)
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting OTEL Collector")
			panic.Go(func() {
				err := otelService.Run(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("Failed to run OTEL Collector")
				}
				_ = in.Shutdowner.Shutdown()
			})
			return nil
		},
		OnStop: func(context.Context) error {
			log.Info().Msg("Stopping OTEL Collector")
			otelService.Shutdown()
			return nil
		},
	})

	return otelService, nil
}
