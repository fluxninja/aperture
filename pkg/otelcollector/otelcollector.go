package otelcollector

import (
	"context"
	"errors"
	"fmt"
	"path"

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

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	"github.com/fluxninja/aperture/pkg/notifiers"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	otelcustom "github.com/fluxninja/aperture/pkg/otelcollector/custom"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
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
	EtcdClient       *etcdclient.Client
	AgentInfo        *agentinfo.AgentInfo
	Factories        otelcol.Factories
	Lifecycle        fx.Lifecycle
	Shutdowner       fx.Shutdowner
	Unmarshaller     config.Unmarshaller
	StatusRegistry   status.Registry
	BaseConfig       *otelconfig.OTelConfig `name:"base"`
	Logger           *log.Logger
	Readiness        *jobs.MultiJob           `name:"readiness.service"`
	ExtensionConfigs []*otelconfig.OTelConfig `group:"extension-config"`
}

// setup creates and runs a new instance of OTel Collector with the passed configuration.
func setup(in ConstructorIn) error {
	policyConfigUnmarshaler := otelconfig.NewOTelConfigUnmarshaler("policy", map[string]interface{}{})
	infraMeters := map[string]*policylangv1.InfraMeter{}
	handleInfraMeterUpdate := func(event notifiers.Event, unmarshaller config.Unmarshaller) {
		log.Info().Str("event", event.String()).Msg("infra meter update")
		infraMeter := &policylangv1.InfraMeter{}
		if err := unmarshaller.UnmarshalKey("", infraMeter); err != nil {
			log.Error().Err(err).Msg("unmarshalling infra meter")
			return
		}
		switch event.Type {
		case notifiers.Write:
			infraMeters[string(event.Key)] = infraMeter
		case notifiers.Remove:
			delete(infraMeters, string(event.Key))
		}
		otelCfg := otelconfig.NewOTelConfig()
		if err := otelcustom.AddCustomMetricsPipelines(otelCfg, infraMeters); err != nil {
			log.Error().Err(err).Msg("unable to add custom metrics pipelines")
			return
		}
		// trigger update
		log.Info().Msgf("received infra meter update, hot re-loading OTel, total infra meters: %d", len(infraMeters))
		policyConfigUnmarshaler.UpdateMap(otelCfg.AsMap())
	}

	// Get Agent Group from host info gatherer
	agentGroupName := in.AgentInfo.GetAgentGroup()
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.TelemetryCollectorConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(in.EtcdClient, etcdPath)
	if err != nil {
		return err
	}
	unmarshalNotifier, err := notifiers.NewUnmarshalPrefixNotifier("",
		handleInfraMeterUpdate,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller)
	if err != nil {
		return fmt.Errorf("creating unmarshal notifier: %w", err)
	}
	notifiers.WatcherLifecycle(in.Lifecycle, watcher, []notifiers.PrefixNotifier{unmarshalNotifier})

	var otelService *otelcol.Collector
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			in.BaseConfig.AddProcessor(otelconsts.ProcessorCustomMetrics, map[string]any{
				"attributes": []map[string]interface{}{
					{
						"key":    "service.name",
						"action": "upsert",
						"value":  "aperture-custom-metrics",
					},
				},
			})

			uris := []string{"service:main", "policy:telemetry-collector"}
			providers := map[string]confmap.Provider{
				"service": otelconfig.NewOTelConfigUnmarshaler("service", in.BaseConfig.AsMap()),
				"policy":  policyConfigUnmarshaler,
			}
			for i, extensionConfig := range in.ExtensionConfigs {
				if extensionConfig == nil {
					continue
				}
				scheme := fmt.Sprintf("extension-%v", i)
				uris = append(uris, fmt.Sprintf("%v:%v", scheme, scheme))
				providers[scheme] = otelconfig.NewOTelConfigUnmarshaler(scheme, extensionConfig.AsMap())
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
			panichandler.Go(func() {
				err := otelService.Run(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("Failed to run OTel Collector")
				}
				_ = in.Shutdowner.Shutdown()
			})
			return nil
		},
		OnStop: func(context.Context) error {
			setReadinessStatus(in.StatusRegistry, nil, errors.New("OTel collector stopping"))
			log.Info().Msg("Stopping OTel Collector")
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
		err = errors.New("otel-collector is not running")
	}
	return wrapperspb.String(state.String()), err
}

func setReadinessStatus(statusRegistry status.Registry, msg proto.Message, err error) {
	statusRegistry.Child("system", "readiness").
		Child("component", "otel-collector").
		SetStatus(status.NewStatus(msg, err))
}
