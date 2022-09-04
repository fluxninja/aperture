package otelcollector

// TODO: organize the constants by their usage.
// example:
// aperture.* are used in ext_authz.CheckResponse.DynamicMetadata
// others are being used to get attributes from traces and logs

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

	// MarshalledAuthzResponseLabel contains JSON encoded response from authz.
	MarshalledAuthzResponseLabel = "aperture.authz_response"

	// AuthzStatusLabel describes the status reported from authz processing.
	AuthzStatusLabel = "authz_status"

	// MarshalledCheckResponseLabel contains JSON encoded check response struct.
	MarshalledCheckResponseLabel = "aperture.check_response"

	// MissingAttributeSourceValue is a special attribute value, which can
	// happen when (eg. Envoy's) logger tries to send attribute value, but its
	// source is missing. Eg. In case authz couldn't reach agent, so we know
	// nothing about flowcontrol policies.  Note that this is different case
	// from "just empty", eg. "", "[]" or "{}".
	MissingAttributeSourceValue = "-"

	// MarshalledLabelsLabel describes labels relevant to this traffic.
	// This is JSON encoded field:
	// {
	//   "foo": "bar",
	//   "fizz": "buzz"
	// }.
	MarshalledLabelsLabel = "aperture.labels"
	// LabeledLabel describes if there are any labels matched to traffic.
	LabeledLabel = "labeled"
	// StatusCodeLabel describes HTTP status code of the response.
	StatusCodeLabel = "http.status_code"
	// DurationLabel describes duration of the HTTP request in milliseconds.
	DurationLabel = "duration_millis"
	// HTTPRequestContentLength describes length of the HTTP request content in bytes.
	HTTPRequestContentLength = "http.request_content_length"
	// HTTPResponseContentLength describes length of the HTTP response content in bytes.
	HTTPResponseContentLength = "http.response_content_length"
	// DecisionTypeLabel describes the decision type taken by policy.
	DecisionTypeLabel = "decision_type"
	// DecisionErrorReasonLabel describes the error reason of the decision taken by policy.
	DecisionErrorReasonLabel = "decision_error_reason"
	// DecisionRejectReasonLabel describes the reject reason of the decision taken by policy.
	DecisionRejectReasonLabel = "decision_reject_reason"
	// RateLimitersLabel describes rate limiters matched to the traffic.
	RateLimitersLabel = "rate_limiters"
	// DroppingRateLimitersLabel describes rate limiters dropping the traffic.
	DroppingRateLimitersLabel = "dropping_rate_limiters"
	// ConcurrencyLimitersLabel describes rate limiters matched to the traffic.
	ConcurrencyLimitersLabel = "concurrency_limiters"
	// DroppingConcurrencyLimitersLabel describes rate limiters dropping the traffic.
	DroppingConcurrencyLimitersLabel = "dropping_concurrency_limiters"
	// FluxMetersLabel describes flux meters matched to the traffic.
	FluxMetersLabel = "flux_meters"
	// FlowLabelKeysLabel describes keys of flow labels matched to the traffic.
	FlowLabelKeysLabel = "flow_label_keys"
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
	// FeatureIDLabel describes the ID of the feature.
	FeatureIDLabel = "feature.id"
	// FeatureStatusLabel describes the status of the feature.
	FeatureStatusLabel = "feature.status"
	// EntityNameLabel describes entity name e.g. pod name.
	EntityNameLabel = "entity_name"
	// TimestampLabel describes timestamp of the request.
	TimestampLabel = "timestamp"
	// AgentGroupLabel describes cluster to which metrics refer.
	AgentGroupLabel = "agent_group"
	// ServicesLabel describes services to which metrics refer.
	ServicesLabel = "services"
)
