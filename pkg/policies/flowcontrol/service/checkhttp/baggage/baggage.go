package baggage

import (
	"fmt"
	"net/url"

	otel_baggage "go.opentelemetry.io/otel/baggage"

	"github.com/fluxninja/aperture/pkg/log"
	flowlabel "github.com/fluxninja/aperture/pkg/policies/flowcontrol/label"
)

// Headers is a header map.
type Headers map[string]string

// Propagator defines how to extract flow labels (baggage) from and put into
// headers.
type Propagator interface {
	// Extract extracts flow labels from headers
	Extract(headers Headers) flowlabel.FlowLabels

	// Inject emits instructions for envoy how to inject flow labels into
	// headers based on given flow labels and existing headers
	//
	// The returned list is expected to be put in
	// CheckResponse.OkHttpResponse.Headers, so that envoy will take care of
	// injecting appropriate headers.
	Inject(flowLabels flowlabel.FlowLabels, headers Headers) (map[string]string, error)
}

// W3Baggage handles baggage propagation in a single `baggage` header, as
// described in https://www.w3.org/TR/baggage/
//
// All baggage items are mapped to flow labels 1:1. This could be tweaked in future:
// * we can use some prefixing/filtering,
// * alternatively, we could put _all_ flow labels as a _single_ baggage item (eg. Fn-flow).
type W3Baggage struct{}

const (
	w3BaggageHeaderName = "baggage"
)

// Extract extracts flow labels from w3Baggage headers.
func (b W3Baggage) Extract(headers Headers) flowlabel.FlowLabels {
	baggage, err := otel_baggage.Parse(headers[w3BaggageHeaderName])
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse baggage header")
		return flowlabel.FlowLabels{}
	}
	flowLabels := make(flowlabel.FlowLabels)
	for _, member := range baggage.Members() {
		value, err := url.QueryUnescape(member.Value())
		if err != nil {
			log.Warn().Msg("Could not unescape flow label value in baggage")
			continue
		}
		flowLabels[member.Key()] = flowlabel.FlowLabelValue{
			Value:     value,
			Telemetry: true,
		}
	}
	return flowLabels
}

// Inject emits instructions for envoy how to inject flow labels into headers supported by baggage propagator.
func (b W3Baggage) Inject(
	flowLabels flowlabel.FlowLabels,
	headers Headers,
) (map[string]string, error) {
	members := make([]otel_baggage.Member, 0, len(flowLabels))
	for k, v := range flowLabels {
		if !v.Telemetry {
			continue
		}
		member, err := otel_baggage.NewMember(k, v.Value)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if len(members) == 0 {
		return nil, nil
	}
	baggage, err := otel_baggage.New(members...)
	if err != nil {
		return nil, err
	}

	baggageValue, baggageAlreadyExists := headers[w3BaggageHeaderName]

	if baggageAlreadyExists {
		baggageValue = fmt.Sprintf("%s,%s", baggageValue, baggage.String())
	} else {
		baggageValue = baggage.String()
	}

	return map[string]string{w3BaggageHeaderName: baggageValue}, nil
}
