package validator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/log"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// TraceHandler implements ExportTraceService.
type TraceHandler struct {
	tracev1.UnimplementedTraceServiceServer
}

// Export is a dummy Export handler.
func (t TraceHandler) Export(ctx context.Context, req *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	log.Trace().Msg("Received Export request")

	var merr error

	for _, resourceSpans := range req.ResourceSpans {
		for _, scopeSpan := range resourceSpans.ScopeSpans {
			for _, span := range scopeSpan.Spans {
				var err error
				var flowStartTS, flowEndTS, workloadTS int64
				for _, attribute := range span.Attributes {
					switch attribute.Key {
					case otelconsts.ApertureCheckResponseLabel:
						log.Trace().Str("attribute", otelconsts.ApertureCheckResponseLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						checkResponse := &flowcontrolv1.CheckResponse{}
						var wireMsg []byte
						if !strings.HasPrefix(v, "{") {
							wireMsg, err = base64.StdEncoding.DecodeString(v)
							if err != nil {
								log.Error().Err(err).Msg("Failed to unmarshal as base64")
							}
							perr := proto.Unmarshal(wireMsg, checkResponse)
							if err != nil {
								log.Error().Err(err).Msg("Failed to unmarshal as protobuf")
								err = multierr.Append(err, fmt.Errorf("invalid %s: %w", otelconsts.ApertureCheckResponseLabel, perr))
							}
							continue
						}
						perr := json.Unmarshal([]byte(v), checkResponse)
						if perr != nil {
							log.Error().Err(err).Msg("Failed to unmarshal as json")
							err = multierr.Append(err, fmt.Errorf("invalid %s: %w", otelconsts.ApertureCheckResponseLabel, perr))
						}
					case otelconsts.ApertureSourceLabel:
						log.Trace().Str("attribute", otelconsts.ApertureSourceLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						if v != "sdk" {
							log.Error().Msg("Failed to validate source")
							err = multierr.Append(err, fmt.Errorf("invalid %s", otelconsts.ApertureSourceLabel))
						}
					case otelconsts.ApertureFlowStatusLabel:
						log.Trace().Str("attribute", otelconsts.ApertureFlowStatusLabel).Msg("Validating attribute")
						v := attribute.Value.GetStringValue()
						if v != otelconsts.ApertureFlowStatusOK && v != otelconsts.ApertureFlowStatusError && v != "Unset" {
							log.Error().Msg("Failed to validate flow status")
							err = multierr.Append(err, fmt.Errorf("invalid %s", otelconsts.ApertureFlowStatusLabel))
						}
					case otelconsts.ApertureFlowStartTimestampLabel:
						flowStartTS = attribute.Value.GetIntValue()
					case otelconsts.ApertureFlowEndTimestampLabel:
						flowEndTS = attribute.Value.GetIntValue()
					case otelconsts.ApertureWorkloadStartTimestampLabel:
						workloadTS = attribute.Value.GetIntValue()
					}
				}
				log.Trace().Str("attribute", otelconsts.ApertureFlowStartTimestampLabel).Int64("value", flowStartTS).Msg("Validating attribute")
				if flowStartTS == 0 {
					log.Error().Msg("Missing start flow timestamp")
					err = multierr.Append(err, fmt.Errorf("invalid %s", otelconsts.ApertureFlowStartTimestampLabel))
				}
				log.Trace().Str("attribute", otelconsts.ApertureFlowEndTimestampLabel).Int64("value", flowEndTS).Msg("Validating attribute")
				if flowEndTS == 0 {
					log.Error().Msg("Missing end flow timestamp")
					err = multierr.Append(err, fmt.Errorf("invalid %s", otelconsts.ApertureFlowEndTimestampLabel))
				}
				if flowStartTS > flowEndTS {
					log.Error().Msg("Failed to validate start and end flow timestamps")
					err = multierr.Append(err, fmt.Errorf("invalid %s and %s", otelconsts.ApertureFlowStartTimestampLabel, otelconsts.ApertureFlowEndTimestampLabel))
				}
				log.Trace().Str("attribute", otelconsts.ApertureWorkloadStartTimestampLabel).Msg("Validating attribute")
				if workloadTS == 0 {
					log.Error().Msg("Failed to validate workload start timestamp")
					err = multierr.Append(err, fmt.Errorf("invalid %s", otelconsts.ApertureWorkloadStartTimestampLabel))
				}
				merr = multierr.Append(merr, err)
				if merr != nil {
					return &tracev1.ExportTraceServiceResponse{}, merr
				}
			}
		}
	}

	log.Info().Msg("Validated span attributes")
	return &tracev1.ExportTraceServiceResponse{}, merr
}
