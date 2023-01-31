package enrichmentprocessor

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// entityNameLabels is a slice of labels using which this processor will try
// to enrich i.e. it will try using all those labels one by one as keys
// in the entity cache.
var entityNameLabels = []string{otelconsts.PodNameLabel}

type enrichmentProcessor struct {
	cache *entitycache.EntityCache
}

func newProcessor(cache *entitycache.EntityCache) *enrichmentProcessor {
	return &enrichmentProcessor{
		cache: cache,
	}
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (ep *enrichmentProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

// Start indicates and logs the start of the enrichment processor.
func (ep *enrichmentProcessor) Start(_ context.Context, _ component.Host) error {
	log.Debug().Msg("enrichment processor start")
	return nil
}

// Shutdown indicates and logs the shutdown of the enrichment processor.
func (ep *enrichmentProcessor) Shutdown(context.Context) error {
	log.Debug().Msg("enrichment processor shutdown")
	return nil
}

// ConsumeMetrics receives pmetric.Metrics for consumption then returns the enriched metrics.
func (ep *enrichmentProcessor) ConsumeMetrics(ctx context.Context, origMd pmetric.Metrics) (pmetric.Metrics, error) {
	if ep.cache == nil {
		return pmetric.Metrics{}, fmt.Errorf("cache not provided")
	}
	md := pmetric.NewMetrics()
	origMd.CopyTo(md)
	// Enrich common attributes in resource
	otelcollector.IterateResourceMetrics(md, func(resourceMetrics pmetric.ResourceMetrics) {
		ep.enrichMetrics(resourceMetrics.Resource().Attributes())
	})
	// Enrich attributes in each of the metric
	otelcollector.IterateMetrics(md, func(metric pmetric.Metric) {
		otelcollector.IterateDataPoints(metric, ep.enrichMetrics)
	})
	return md, nil
}

func (ep *enrichmentProcessor) enrichMetrics(attributes pcommon.Map) {
	for _, label := range entityNameLabels {
		hostNamex, ok := attributes.Get(label)
		if !ok {
			continue
		}
		hostName := hostNamex.Str()
		hostEntity, err := ep.cache.GetByName(hostName)
		attributes.Remove(label)
		if err != nil {
			log.Trace().Str("label", label).Str("name", hostName).Msg("Skipping because entity not found in cache")
			continue
		}
		// We don't want this to be OTLP slice, as it is weirdly formatten when written to prometheus.
		attributes.PutStr(otelconsts.ApertureServicesLabel, strings.Join(hostEntity.Services, ","))
	}
}
