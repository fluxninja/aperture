package validator

import (
	"context"
	"fmt"

	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// TraceHandler implements ExportTraceService.
type TraceHandler struct {
	tracev1.UnimplementedTraceServiceServer

	Rejects  int64
	Rejected int64
}

// Export is a dummy Export handler.
func (t TraceHandler) Export(ctx context.Context, req *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	log.Info().Msg("Received Export request")

	for _, resourceSpans := range req.ResourceSpans {
		for _, scopeSpan := range resourceSpans.ScopeSpans {
			for _, span := range scopeSpan.Spans {
				for _, attribute := range span.Attributes {
					switch attribute.Key {
					case otelcollector.ApertureCheckResponseLabel:
						v := attribute.Value.GetStringValue()
						checkResponse := &flowcontrolv1.CheckResponse{}
						err := protojson.Unmarshal([]byte(v), checkResponse)
						if err != nil {
							log.Error().Err(err).Msg("Failed to validate flowcontrol CheckResponse")
							return &tracev1.ExportTraceServiceResponse{}, fmt.Errorf("invalid %s: %w", otelcollector.ApertureCheckResponseLabel, err)
						}
					case otelcollector.ApertureControlPointLabel:
						v := attribute.Value.GetStringValue()
						if v != "sdk" {
							log.Error().Msg("Failed to validate control point")
							return &tracev1.ExportTraceServiceResponse{}, fmt.Errorf("invalid %s", otelcollector.ApertureControlPointLabel)
						}
					}
				}
			}
		}
	}

	return &tracev1.ExportTraceServiceResponse{}, nil
}
