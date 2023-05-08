package agent

import (
	"crypto/tls"
	"fmt"
	"path"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/fx"
	"k8s.io/client-go/rest"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	agentconfig "github.com/fluxninja/aperture/cmd/aperture-agent/config"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/notifiers"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	inframeter "github.com/fluxninja/aperture/pkg/otelcollector/infra-meter"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/utils"
)

func provideAgent(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
	ai *agentinfo.AgentInfo,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
	shutdowner fx.Shutdowner,
) (*otelconfig.OTelConfigProvider, *otelconfig.OTelConfigProvider, error) {
	var agentCfg agentconfig.AgentOTelConfig
	if err := unmarshaller.UnmarshalKey("otel", &agentCfg); err != nil {
		return nil, nil, fmt.Errorf("unmarshalling otel config: %w", err)
	}

	otelCfg := otelconfig.NewOTelConfig()
	otelCfg.SetDebugPort(&agentCfg.CommonOTelConfig)
	otelCfg.AddDebugExtensions(&agentCfg.CommonOTelConfig)

	addLogsPipeline(otelCfg, &agentCfg)
	addTracesPipeline(otelCfg, lis)
	addMetricsPipeline(otelCfg, &agentCfg, tlsConfig, lis, promClient)

	customConfig := map[string]*policylangv1.InfraMeter{}
	if !agentCfg.DisableKubeletScraper {
		customConfig[otelconsts.ReceiverKubeletStats] = inframeter.InfraMeterForKubeletStats()
	}

	if err := inframeter.AddInfraMeters(otelCfg, customConfig); err != nil {
		return nil, nil, fmt.Errorf("adding custom metrics pipelines: %w", err)
	}
	otelconfig.AddAlertsPipeline(otelCfg, agentCfg.CommonOTelConfig, otelconsts.ProcessorAgentResourceLabels)

	baseConfigProvider := otelconfig.NewOTelConfigProvider("service", otelCfg)

	tcConfigProvider := otelconfig.NewOTelConfigProvider("telemetry-collector", otelconfig.NewOTelConfig())

	allInfraMeters := map[string]map[string]*policylangv1.InfraMeter{}
	handleInfraMeterUpdate := func(event notifiers.Event, unmarshaller config.Unmarshaller) {
		log.Info().Str("event", event.String()).Msg("infra meter update")
		tc := &policylangv1.TelemetryCollector{}
		if err := unmarshaller.UnmarshalKey("", tc); err != nil {
			log.Error().Err(err).Msg("unmarshalling telemetry collector")
			return
		}
		infraMeters := tc.GetInfraMeters()
		key := string(event.Key)

		switch event.Type {
		case notifiers.Write:
			allInfraMeters[key] = infraMeters
		case notifiers.Remove:
			delete(allInfraMeters, key)
		}
		otelCfg := otelconfig.NewOTelConfig()
		ims := map[string]*policylangv1.InfraMeter{}
		for prefix, v := range allInfraMeters {
			for k, v := range v {
				ims[fmt.Sprintf("%s/%s", prefix, k)] = v
			}
		}
		if err := inframeter.AddInfraMeters(otelCfg, ims); err != nil {
			log.Error().Err(err).Msg("unable to add custom metrics pipelines")
			utils.Shutdown(shutdowner)
			return
		}
		// trigger update
		log.Info().Msgf("received infra meter update, hot re-loading OTel, total infra meters: %d", len(allInfraMeters))
		tcConfigProvider.UpdateConfig(otelCfg)
	}

	// Get Agent Group from host info gatherer
	agentGroupName := ai.GetAgentGroup()
	// Scope the sync to the agent group.
	etcdPath := path.Join(paths.TelemetryCollectorConfigPath,
		paths.AgentGroupPrefix(agentGroupName))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, nil, err
	}
	unmarshalNotifier, err := notifiers.NewUnmarshalPrefixNotifier("",
		handleInfraMeterUpdate,
		config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller)
	if err != nil {
		return nil, nil, fmt.Errorf("creating unmarshal notifier: %w", err)
	}
	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{unmarshalNotifier})

	return baseConfigProvider, tcConfigProvider, nil
}

func addLogsPipeline(
	config *otelconfig.OTelConfig,
	userConfig *agentconfig.AgentOTelConfig,
) {
	// Common dependencies for pipelines
	config.AddReceiver(otelconsts.ReceiverOTLP, otlpreceiver.Config{})
	// Note: Passing map[string]interface{}{} instead of real config, so that
	// processors' configs' default work.
	config.AddProcessor(otelconsts.ProcessorMetrics, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPrerollup,
		userConfig.BatchPrerollup.Timeout.AsDuration(),
		userConfig.BatchPrerollup.SendBatchSize,
		userConfig.BatchPrerollup.SendBatchMaxSize,
	)
	config.AddProcessor(otelconsts.ProcessorRollup, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPostrollup,
		userConfig.BatchPostrollup.Timeout.AsDuration(),
		userConfig.BatchPostrollup.SendBatchSize,
		userConfig.BatchPostrollup.SendBatchMaxSize,
	)
	config.AddExporter(otelconsts.ExporterLogging, nil)

	processors := []string{
		otelconsts.ProcessorMetrics,
		otelconsts.ProcessorBatchPrerollup,
		otelconsts.ProcessorRollup,
		otelconsts.ProcessorBatchPostrollup,
		otelconsts.ProcessorAgentGroup,
	}

	config.Service.AddPipeline("logs", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverOTLP, otelconsts.ConnectorAdapter},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging},
	})
}

func addTracesPipeline(config *otelconfig.OTelConfig, _ *listener.Listener) {
	config.AddConnector(otelconsts.ConnectorAdapter, map[string]any{})
	config.Service.AddPipeline("traces", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverOTLP},
		Exporters: []string{otelconsts.ConnectorAdapter},
	})
}

func addMetricsPipeline(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, agentConfig, tlsConfig, lis)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/fast", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-self", tlsConfig, lis),
		otelconfig.BuildOTelScrapeConfig("aperture-otel", agentConfig.CommonOTelConfig),
	}

	if !agentConfig.DisableKubernetesScraper {
		_, err := rest.InClusterConfig()
		if err == rest.ErrNotInCluster {
			log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
		} else if err != nil {
			log.Warn().Err(err).Msg("Error when discovering k8s environment")
		} else {
			log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
			scrapeConfigs = append(scrapeConfigs, buildKubernetesPodsScrapeConfig())
		}
	}

	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTel. Need to use bare maps instead.
	config.AddReceiver(otelconsts.ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "1s",
				"scrape_timeout":      "1s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func buildKubernetesPodsScrapeConfig() map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-pods",
		"scheme":       "http",
		"metrics_path": "/metrics",
		"kubernetes_sd_configs": []map[string]any{
			{"role": "pod"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_pod_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
			// Scrape only pods which have github.com/fluxninja/scrape=true annotation.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_aperture_tech_scrape"},
				"action":        "keep",
				"regex":         "true",
			},
			// Allow rewrite of scheme, path and port where prometheus metrics are served.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
				"action":        "replace",
				"regex":         "(https?)",
				"target_label":  "__scheme__",
			},
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
				"action":        "replace",
				"target_label":  "__metrics_path__",
				"regex":         "(.+)",
			},
			{
				"source_labels": []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
				"action":        "replace",
				"regex":         `([^:]+)(?::\d+)?;(\d+)`,
				"replacement":   "$$1:$$2",
				"target_label":  "__address__",
			},
		},
		"metric_relabel_configs": []map[string]any{
			// For now, dropping everything. In future, we'll want to filter in some
			// metrics based on policies. See #4632.
			{
				"source_labels": []string{},
				"action":        "drop",
			},
		},
	}
}
