package otelcollector

// TODO (hasit): These are commong envoy keys that need to be moved to a common package

const (
	// ControlPointLabel describes control point which reported traffic.
	// May be 'ingress', 'egress' or 'feature'.
	ControlPointLabel = "control_point"
	// ControlPointIngress const for ingress control point.
	ControlPointIngress = "ingress"
	// ControlPointEgress const for egress control point.
	ControlPointEgress = "egress"
	// ControlPointFeature const for feature control point.
	ControlPointFeature = "feature"

	// LimiterDecisionsLabel describes policies relevant to this traffic.
	// This is JSON encoded field:
	// [
	//   {
	//     "id": "<id of the policy>",
	//     "dropped": true,
	//     "workload": "foo:bar",
	//   },
	//   ...
	// ].
	//
	LimiterDecisionsLabel = "aperture.limiter_decisions"
	// FluxMetersLabel describes the flux meter IDs matched to this traffic.
	FluxMetersLabel = "aperture.fluxmeters"

	// MissingAttributeSourceValue is a special attribute value, which can
	// happen when (eg. Envoy's) logger tries to send attribute value, but its
	// source is missing. Eg. In case authz couldn't reach agent, so we know
	// nothing about flowcontrol policies.  Note that this is different case
	// from "just empty", eg. "", "[]" or "{}".
	MissingAttributeSourceValue = "-"

	// FlowLabel describes flow labels relevant to this traffic.
	// This is JSON encoded field:
	// {
	//   "foo": "bar",
	//   "fizz": "buzz"
	// }.
	//
	FlowLabel = "aperture.flow"
	// LabeledLabel describes if there are any flow labels matched to traffic.
	LabeledLabel = "labeled"
	// StatusCodeLabel describes HTTP status code of the response.
	StatusCodeLabel = "http.status_code"
	// HTTPDurationLabel describes duration of the HTTP request in milliseconds.
	HTTPDurationLabel = "http.duration_millis"
	// HTTPRequestContentLength describes length of the HTTP request content in bytes.
	HTTPRequestContentLength = "http.request_content_length"
	// HTTPResponseContentLength describes length of the HTTP response content in bytes.
	HTTPResponseContentLength = "http.response_content_length"
	// FeatureDurationLabel describes duration of the feature in milliseconds.
	FeatureDurationLabel = "feature.duration_millis"
	// PoliciesMatchedLabel describes if there are any policies matched to traffic.
	PoliciesMatchedLabel = "policy_ids_matched"
	// PoliciesDroppedLabel describes if there are any policies dropped to traffic.
	PoliciesDroppedLabel = "policy_ids_dropped"
	// HostAddressLabel describes host address of the request.
	HostAddressLabel = "net.host.address"
	// PeerAddressLabel describes peer address of the request.
	PeerAddressLabel = "net.peer.address"
	// HostIPLabel describes host IP address of the request.
	HostIPLabel = "net.host.ip"
	// PeerIPLabel describes peer IP address of the request.
	PeerIPLabel = "net.peer.ip"
	// FeatureAddressLabel describes feature address of the request.
	FeatureAddressLabel = "feature.ip"
	// EntityNameLabel describes entity name e.g. pod name.
	EntityNameLabel = "entity_name"
	// TimestampLabel describes timestamp of the request.
	TimestampLabel = "timestamp"
	// AgentGroupLabel describes cluster to which metrics refer.
	AgentGroupLabel = "agent_group"
	// NamespaceLabel describes namespace to which metrics refer.
	NamespaceLabel = "namespace"
	// ServicesLabel describes services to which metrics refer. This is comma-separated
	// list.
	ServicesLabel = "services"
)
