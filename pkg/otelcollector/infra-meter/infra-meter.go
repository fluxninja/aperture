package inframeter

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/leaderonlyreceiver"
	"go.opentelemetry.io/collector/component"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/structpb"
)

// AddInfraMeters adds infra metrics pipelines to the given OTelConfig.
func AddInfraMeters(
	config *otelconfig.Config,
	infraMeters map[string]*policysyncv1.InfraMeterWrapper,
) error {
	if infraMeters == nil {
		infraMeters = map[string]*policysyncv1.InfraMeterWrapper{}
	}
	for pipelineName, infraMeterWrapper := range infraMeters {
		infraMeter := infraMeterWrapper.GetInfraMeter()
		if err := addInfraMeter(
			config, infraMeterWrapper.GetPolicyName(), pipelineName, infraMeterWrapper.GetInfraMeterName(), infraMeter); err != nil {
			return fmt.Errorf("failed to add infra metric pipeline %s: %w", pipelineName, err)
		}
	}
	return nil
}

func addInfraMeter(
	config *otelconfig.Config,
	policyName string,
	pipelineName string,
	infraMeterName string,
	infraMeter *policylangv1.InfraMeter,
) error {
	processorName := fmt.Sprintf("%s-%s-%s", otelconsts.ProcessorInfraMeter, policyName, infraMeterName)
	config.AddProcessor(processorName, map[string]any{
		"attributes": []map[string]interface{}{
			{
				"key":    "service.name",
				"action": "upsert",
				"value":  "aperture-infra-meter",
			},
			{
				"key":    metrics.InfraMeterNameLabel,
				"action": "upsert",
				"value":  infraMeterName,
			},
			{
				"key":    metrics.PolicyNameLabel,
				"action": "upsert",
				"value":  policyName,
			},
		},
	})
	pipelineName = strings.TrimPrefix(pipelineName, "metrics/")

	receiverIDs := map[string]string{}
	processorIDs := map[string]string{}

	for origName, receiverConfig := range infraMeter.Receivers {
		var id component.ID
		if err := id.UnmarshalText([]byte(origName)); err != nil {
			return fmt.Errorf("invalid id %q: %w", origName, err)
		}

		sum := sha256.Sum256([]byte(receiverConfig.String()))
		id = component.NewIDWithName(id.Type(), fmt.Sprintf("%x", sum))
		strID := id.String()

		// If receiver is already present with given id and per-agent-group = false, skip adding receiver with per-agent-group = true.
		if _, ok := config.Receivers[strID]; ok && infraMeter.PerAgentGroup {
			receiverIDs[origName] = strID
			continue
		}

		var cfg any
		cfg = receiverConfig.AsMap()
		id, cfg = leaderonlyreceiver.WrapConfigIf(infraMeter.PerAgentGroup, id, cfg)
		strID = id.String()

		// Remove receiver for per-agent-group infra-meter if a receiver with the same config already exists without per-agent-group.
		if !infraMeter.PerAgentGroup {
			updatedID, _ := leaderonlyreceiver.WrapConfig(id, cfg)
			if _, ok := config.Receivers[updatedID.String()]; ok {
				delete(config.Receivers, updatedID.String())
				for key, value := range receiverIDs {
					if value == updatedID.String() {
						receiverIDs[key] = id.String()
					}
				}

				// Update pipelines to use the new receiver id.
				for key, value := range config.Service.Pipelines {
					if slices.Contains(value.Receivers, updatedID.String()) {
						value.Receivers[slices.Index(value.Receivers, updatedID.String())] = id.String()
						config.Service.Pipelines[key] = value
					}
				}
			}
		}
		receiverIDs[origName] = strID
		config.AddReceiver(strID, cfg)
	}

	for origName, processorConfig := range infraMeter.Processors {
		var id component.ID
		if err := id.UnmarshalText([]byte(origName)); err != nil {
			return fmt.Errorf("invalid id %q: %w", origName, err)
		}

		selectorsList := []interface{}{}
		var selectors *structpb.Value
		if id.Type() == otelconsts.ProcessorK8sAttributes {
			var ok bool
			selectors, ok = processorConfig.Fields[otelconsts.ProcessorK8sAttributesSelectors]
			if ok && selectors != nil {
				selectorsList = selectors.GetListValue().AsSlice()
				delete(processorConfig.Fields, otelconsts.ProcessorK8sAttributesSelectors)
			}
		}

		sum := sha256.Sum256([]byte(processorConfig.String()))
		id = component.NewIDWithName(id.Type(), fmt.Sprintf("%x", sum))
		strID := id.String()
		if processor, ok := config.Processors[strID]; ok {
			if id.Type() == otelconsts.ProcessorK8sAttributes {
				processorConfig.Fields[otelconsts.ProcessorK8sAttributesSelectors] = selectors
				updateK8sAttributesProcessor(processor, selectorsList)
				config.Processors[strID] = processor
			}
			processorIDs[origName] = strID
			continue
		} else {
			if id.Type() == otelconsts.ProcessorK8sAttributes {
				processorConfig.Fields[otelconsts.ProcessorK8sAttributesSelectors] = selectors
			}
			processorIDs[origName] = strID
			var cfg any = processorConfig.AsMap()
			config.AddProcessor(strID, cfg)
		}
	}

	if infraMeter.Pipeline == nil {
		// We treat empty pipeline the same way as not-set pipeline, normalize.
		// This also allows to avoid nil checks below.
		infraMeter.Pipeline = &policylangv1.InfraMeter_MetricsPipeline{}
	}

	if len(infraMeter.Pipeline.Receivers) == 0 && len(infraMeter.Pipeline.Processors) == 0 {
		if len(infraMeter.Processors) >= 1 {
			return fmt.Errorf("empty pipeline, inferring pipeline is supported only with 0 or 1 processors")
		}

		// Skip adding pipeline if there are no receivers and processors.
		if len(infraMeter.Receivers) == 0 && len(infraMeter.Processors) == 0 {
			return nil
		}

		// When pipeline not set explicitly, create pipeline with all defined receivers and processors.
		if len(infraMeter.Receivers) > 0 {
			infraMeter.Pipeline.Receivers = maps.Keys(infraMeter.Receivers)
			sort.Strings(infraMeter.Pipeline.Receivers)
		}
		if len(infraMeter.Processors) > 0 {
			infraMeter.Pipeline.Processors = maps.Keys(infraMeter.Processors)
		}
	}

	config.Service.AddPipeline(normalizePipelineName(pipelineName), otelconfig.Pipeline{
		Receivers: mapSlice(receiverIDs, infraMeter.Pipeline.Receivers),
		Processors: append(
			mapSlice(processorIDs, infraMeter.Pipeline.Processors),
			processorName,
			otelconsts.ProcessorAgentResourceLabels,
		),
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
	return nil
}

// normalizePipelineName normalizes user defined pipeline name by adding
// `metrics/user-defined-` prefix.
// This ensures no builtin metrics pipeline is overwritten.
func normalizePipelineName(pipelineName string) string {
	return fmt.Sprintf("metrics/user-defined-%s", pipelineName)
}

func mapSlice(mapping map[string]string, xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		y, ok := mapping[x]
		if !ok {
			y = x
		}
		ys = append(ys, y)
	}
	return ys
}

func updateK8sAttributesProcessor(processor interface{}, selectorsList []interface{}) {
	processorMap := processor.(map[string]interface{})
	if processorMap != nil {
		var existingSelectorsList []interface{}
		existingSelectors, ok := processorMap[otelconsts.ProcessorK8sAttributesSelectors]
		if !ok || existingSelectors == nil {
			existingSelectorsList = []interface{}{}
		} else {
			existingSelectorsList = existingSelectors.([]interface{})
			if existingSelectorsList == nil {
				existingSelectorsList = []interface{}{}
			}
		}

		if len(selectorsList) > 0 {
			processorMap[otelconsts.ProcessorK8sAttributesSelectors] = append(existingSelectorsList, selectorsList...)
		}
	}
}
