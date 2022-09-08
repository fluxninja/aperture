package baggage

import (
	"net/url"
	"strings"

	envoy_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	otel_baggage "go.opentelemetry.io/otel/baggage"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/log"
	class "github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/resources/classifier/compiler"
)

// Headers is a header map in authz convention – keys are lowercase.
type Headers map[string]string

// Propagator defines how to extract flow labels (baggage) from and put into
// headers
//
// This interface is similar to go.opentelemetry.io/otel/propagation.TextMapPropagator:
// https://github.com/open-telemetry/opentelemetry-go/blob/v1.2.0/propagation/propagation.go#L91-L111
// but is designed to use authz-compatible types and conventions.
type Propagator interface {
	// Extract extracts flow labels from headers
	//
	// The headers are expected to be in envoy's authz convention – with
	// lower-case keys
	Extract(headers Headers) class.FlowLabels

	// Inject emits instructions for envoy how to inject flow labels into
	// headers based on given flow labels and existing headers
	//
	// The returned list is expected to be put in
	// CheckResponse.OkHttpResponse.Headers, so that envoy will take care of
	// injecting appropriate headers.
	Inject(flowLabels class.FlowLabels, headers Headers) ([]*envoy_core.HeaderValueOption, error)
}

// Prefixed puts each flow label into a separate header.  Header name is
// concatenation of Prefix and flow label name.
type Prefixed struct {
	// Header prefix, eg. "uberctx-" – in lower-case, as per envoy's authz
	// convention – (note that this differs from golang's convetions))
	Prefix string
}

// Extract extracts prefixed flow labels from headers.
func (p Prefixed) Extract(headers Headers) class.FlowLabels {
	flowLabels := make(class.FlowLabels)
	for key, val := range headers {
		if strings.HasPrefix(key, p.Prefix) {
			metaKey := strings.TrimPrefix(key, p.Prefix)
			metaVal, err := url.QueryUnescape(val)
			if err != nil {
				log.Warn().Msg("Could not unescape flow label value in baggage")
			} else {
				flowLabels[metaKey] = class.FlowLabelValue{
					Value: metaVal,
					Flags: compiler.LabelFlags{
						Propagate: true,
					},
				}
			}
		}
	}
	return flowLabels
}

// Inject emits instructions for envoy how to inject flow labels into headers supported by prefixed propagator.
func (p Prefixed) Inject(
	flowLabels class.FlowLabels,
	headers Headers,
) ([]*envoy_core.HeaderValueOption, error) {
	newHeaders := make([]*envoy_core.HeaderValueOption, 0, len(flowLabels))
	for key, fl := range flowLabels {
		if !fl.Flags.Propagate {
			continue
		}
		if fl.Flags.Hidden {
			log.Warn().Msg("Hidden flow labels are not supported by Prefixed propagator")
		}
		baggageKey := p.Prefix + key
		newHeader := &envoy_core.HeaderValueOption{
			Header: &envoy_core.HeaderValue{
				Key: baggageKey,
				// Note: not urlescaping the value – envoy will do it by itself.
				Value: fl.Value,
			},
			Append: wrapperspb.Bool(false),
		}
		newHeaders = append(newHeaders, newHeader)
	}
	return newHeaders, nil
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
	hiddenPropertyKey   = "hidden"
)

// Extract extracts flow labels from w3Baggage headers.
func (b W3Baggage) Extract(headers Headers) class.FlowLabels {
	baggage, err := otel_baggage.Parse(headers[w3BaggageHeaderName])
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse baggage header")
		return nil
	}
	flowLabels := make(class.FlowLabels)
	for _, member := range baggage.Members() {
		value, err := url.QueryUnescape(member.Value())
		if err != nil {
			log.Warn().Msg("Could not unescape flow label value in baggage")
			continue
		}
		isHidden := false
		for _, prop := range member.Properties() {
			if prop.Key() == hiddenPropertyKey {
				isHidden = true
			}
		}
		flowLabels[member.Key()] = class.FlowLabelValue{
			Value: value,
			Flags: compiler.LabelFlags{
				Hidden:    isHidden,
				Propagate: true,
			},
		}
	}
	return flowLabels
}

// Inject emits instructions for envoy how to inject flow labels into headers supported by baggage propagator.
func (b W3Baggage) Inject(
	flowLabels class.FlowLabels,
	headers Headers,
) ([]*envoy_core.HeaderValueOption, error) {
	members := make([]otel_baggage.Member, 0, len(flowLabels))
	for k, v := range flowLabels {
		if !v.Flags.Propagate {
			continue
		}
		var props []otel_baggage.Property
		if v.Flags.Hidden {
			prop, err := otel_baggage.NewKeyProperty(hiddenPropertyKey)
			if err != nil {
				log.Panic().Err(err).Msgf("Failed to create new key property: %v", err)
			}
			props = append(props, prop)
		}
		member, err := otel_baggage.NewMember(k, v.Value, props...)
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
	_, baggageAlreadyExists := headers[w3BaggageHeaderName]
	return []*envoy_core.HeaderValueOption{{
		Header: &envoy_core.HeaderValue{
			Key:   w3BaggageHeaderName,
			Value: baggage.String(),
		},
		Append: wrapperspb.Bool(baggageAlreadyExists),
	}}, nil
}
