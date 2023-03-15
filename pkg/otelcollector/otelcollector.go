package otelcollector

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/otelcol"
	logsv1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	metricsv1 "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/status"
)

// Module is a fx module that invokes OTEL Collector.
func Module() fx.Option {
	return fx.Options(
		grpcgateway.RegisterHandler{Handler: logsv1.RegisterLogsServiceHandlerFromEndpoint}.Annotate(),
		grpcgateway.RegisterHandler{Handler: tracev1.RegisterTraceServiceHandlerFromEndpoint}.Annotate(),
		grpcgateway.RegisterHandler{Handler: metricsv1.RegisterMetricsServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(setup),
	)
}

// ConstructorIn describes parameters passed to create OTEL Collector, server providing the OpenTelemetry Collector service.
type ConstructorIn struct {
	fx.In
	Factories        otelcol.Factories
	Lifecycle        fx.Lifecycle
	Shutdowner       fx.Shutdowner
	Unmarshaller     config.Unmarshaller
	StatusRegistry   status.Registry
	BaseConfig       *otelconfig.OTELConfig `name:"base"`
	Logger           *log.Logger
	Readiness        *jobs.MultiJob           `name:"readiness.service"`
	ExtensionConfigs []*otelconfig.OTELConfig `group:"extension-config"`
}

// setup creates and runs a new instance of OTEL Collector with the passed configuration.
func setup(in ConstructorIn) error {
	uris := []string{"file:main"}
	var otelService *otelcol.Collector
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			setReadinessStatus(in.StatusRegistry, nil, errors.New("OTEL collector starting"))
			providers := map[string]confmap.Provider{
				"file": otelconfig.NewOTELConfigUnmarshaler(in.BaseConfig.AsMap()),
			}
			for i, extensionConfig := range in.ExtensionConfigs {
				if extensionConfig == nil {
					continue
				}
				scheme := fmt.Sprintf("extension-%v", i)
				uris = append(uris, fmt.Sprintf("%v:%v", scheme, scheme))
				providers[scheme] = otelconfig.NewOTELConfigUnmarshaler(extensionConfig.AsMap())
			}

			configProvider, err := otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
				ResolverSettings: confmap.ResolverSettings{
					URIs:      uris,
					Providers: providers,
					Converters: []confmap.Converter{
						expandconverter.New(),
					},
				},
			})
			if err != nil {
				return fmt.Errorf("creating OTEL config provider: %w", err)
			}
			otelService, err = otelcol.NewCollector(
				otelcol.CollectorSettings{
					BuildInfo:               component.NewDefaultBuildInfo(),
					Factories:               in.Factories,
					ConfigProvider:          configProvider,
					DisableGracefulShutdown: true,
					LoggingOptions: []zap.Option{zap.WrapCore(func(zapcore.Core) zapcore.Core {
						return log.NewZapAdapter(in.Logger, "otel-collector")
					})},
					// NOTE: do not remove this becauase it causes a data-race condition.
					SkipSettingGRPCLogger: true,
				},
			)
			if err != nil {
				return fmt.Errorf("constructing OTEL Service: %v", err)
			}
			err = registerReadinessJob(in.StatusRegistry, in.Readiness, otelService)
			if err != nil {
				return fmt.Errorf("registering OTEL Service readiness job: %v", err)
			}

			log.Info().Msg("Starting OTEL Collector")
			panichandler.Go(func() {
				err := otelService.Run(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("Failed to run OTEL Collector")
				}
				_ = in.Shutdowner.Shutdown()
			})
			return nil
		},
		OnStop: func(context.Context) error {
			setReadinessStatus(in.StatusRegistry, nil, errors.New("OTEL collector stopping"))
			log.Info().Msg("Stopping OTEL Collector")
			otelService.Shutdown()
			return nil
		},
	})

	return nil
}

func registerReadinessJob(
	statusRegistry status.Registry,
	readiness *jobs.MultiJob,
	otelService *otelcol.Collector,
) error {
	return readiness.RegisterJob(jobs.NewBasicJob("otel-collector", func(ctx context.Context) (proto.Message, error) {
		msg, err := otelState(otelService)
		setReadinessStatus(statusRegistry, msg, err)
		return msg, err
	}))
}

func otelState(otelService *otelcol.Collector) (proto.Message, error) {
	state := otelService.GetState()
	var err error
	if state != otelcol.StateRunning {
		err = errors.New("otel-collector is not running")
	}
	return wrapperspb.String(state.String()), err
}

func setReadinessStatus(statusRegistry status.Registry, msg proto.Message, err error) {
	statusRegistry.Child("system", "readiness").
		Child("component", "otel-collector").
		SetStatus(status.NewStatus(msg, err))
}
