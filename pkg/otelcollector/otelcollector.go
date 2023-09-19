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

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/net/grpcgateway"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module is a fx module that invokes OTel Collector.
func Module() fx.Option {
	return fx.Options(
		grpcgateway.RegisterHandler{Handler: logsv1.RegisterLogsServiceHandlerFromEndpoint}.Annotate(),
		grpcgateway.RegisterHandler{Handler: tracev1.RegisterTraceServiceHandlerFromEndpoint}.Annotate(),
		grpcgateway.RegisterHandler{Handler: metricsv1.RegisterMetricsServiceHandlerFromEndpoint}.Annotate(),
		fx.Invoke(setup),
	)
}

// ConstructorIn describes parameters passed to create OTel Collector, server providing the OpenTelemetry Collector service.
type ConstructorIn struct {
	fx.In
	Factories      otelcol.Factories
	Lifecycle      fx.Lifecycle
	Shutdowner     fx.Shutdowner
	Unmarshaller   config.Unmarshaller
	StatusRegistry status.Registry
	ConfigProvider *otelconfig.Provider
	Logger         *log.Logger
	Readiness      *jobs.MultiJob `name:"readiness.service"`
}

// setup creates and runs a new instance of OTel Collector with the passed configuration.
func setup(in ConstructorIn) error {
	var otelService *otelcol.Collector
	var wg panichandler.WaitGroup
	var cancelRun context.CancelFunc // See OnStop why we need this.
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			scheme := in.ConfigProvider.Scheme()
			uris := []string{fmt.Sprintf("%v:%v", scheme, scheme)}
			providers := map[string]confmap.Provider{
				scheme: in.ConfigProvider,
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
				return fmt.Errorf("creating OTel config provider: %w", err)
			}

			setReadinessStatus(in.StatusRegistry, nil, errors.New("OTel collector starting"))
			otelService, err = otelcol.NewCollector(
				otelcol.CollectorSettings{
					BuildInfo:               component.NewDefaultBuildInfo(),
					Factories:               in.Factories,
					ConfigProvider:          configProvider,
					DisableGracefulShutdown: true,
					LoggingOptions: []zap.Option{zap.WrapCore(func(zapcore.Core) zapcore.Core {
						return log.NewZapAdapter(in.Logger, "otel-collector")
					})},
					// NOTE: do not remove this because it causes a data-race condition.
					SkipSettingGRPCLogger: true,
				},
			)
			if err != nil {
				return fmt.Errorf("constructing OTel Service: %v", err)
			}
			err = registerReadinessJob(in.StatusRegistry, in.Readiness, otelService)
			if err != nil {
				return fmt.Errorf("registering OTel Service readiness job: %v", err)
			}

			log.Info().Msg("Starting OTel Collector")
			var runCtx context.Context
			runCtx, cancelRun = context.WithCancel(context.Background())
			wg.Go(func() {
				err := otelService.Run(runCtx)
				if err != nil {
					log.Error().Err(err).Msg("Failed to run OTel Collector")
				}
				utils.Shutdown(in.Shutdowner)
			})
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			setReadinessStatus(in.StatusRegistry, nil, errors.New("OTel collector stopping"))
			log.Info().Msg("Stopping OTel Collector")

			var cancelStop context.CancelFunc
			stopCtx, cancelStop = context.WithCancel(stopCtx)
			defer cancelStop()
			go func() {
				<-stopCtx.Done()
				// While collector can be stopped by calling Shutdown(), the
				// context used for actual shutting down is the context passed
				// to collector.Run(). To make sure shutdown is actually
				// bounded by given stopCtx, we need this cancel to enforce
				// stopCtx's deadline.
				cancelRun()
			}()

			otelService.Shutdown()
			// Shutdown only starts the shutdown process, we need to wait for
			// actual shutdown to finish.
			wg.Wait()
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
	return readiness.RegisterJob(jobs.NewBasicJob("otel-collector", func(context.Context) (proto.Message, error) {
		msg, err := otelState(otelService)
		setReadinessStatus(statusRegistry, msg, err)
		return msg, err
	}))
}

func otelState(otelService *otelcol.Collector) (proto.Message, error) {
	state := otelService.GetState()
	var err error
	if state != otelcol.StateRunning {
		err = fmt.Errorf("otel-collector is not running. State: %s", state.String())
	}
	return wrapperspb.String(state.String()), err
}

func setReadinessStatus(statusRegistry status.Registry, msg proto.Message, err error) {
	statusRegistry.Child("system", "readiness").
		Child("component", "otel-collector").
		SetStatus(status.NewStatus(msg, err))
}
