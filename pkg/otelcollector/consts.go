package otelcollector

// TODO: organize the constants by their usage.
// example:
// aperture.* are used in ext_authz.CheckResponse.DynamicMetadata
// others are being used to get attributes from traces and logs

const (
	/* Common labels available on all check calls. */

	// MarshalledLabelsLabel describes labels relevant to this traffic.
	// This is JSON encoded field:
	// {
	//   "foo": "bar",
	//   "fizz": "buzz"
	// }.
	MarshalledLabelsLabel = "aperture.labels"
	// DurationLabel describes duration of the flow in milliseconds.
	// NOTE: not available on spans because the span timestamps are encoded natively.
	DurationLabel = "duration_millis"
	// MarshalledCheckResponseLabel contains JSON encoded check response struct.
	MarshalledCheckResponseLabel = "aperture.check_response"
	// ControlPointLabel describes control point which reported traffic.
	// May be 'ingress', 'egress' or 'feature'.
	ControlPointLabel = "control_point"
	// ControlPointIngress const for ingress control point.
	ControlPointIngress = "ingress"
	// ControlPointEgress const for egress control point.
	ControlPointEgress = "egress"
	// ControlPointFeature const for feature control point.
	ControlPointFeature = "feature"

	/* Derived label that is applied based on content of labels. */

	// LabeledLabel describes if there are any labels matched to traffic.
	LabeledLabel = "labeled"

	/* Envoy specific Authz label. */

	// MarshalledAuthzResponseLabel contains JSON encoded response from authz.
	MarshalledAuthzResponseLabel = "aperture.authz_response"

	/* Derived label based on content of authz response label. */

	// AuthzStatusLabel describes the status reported from authz processing.
	AuthzStatusLabel = "authz_status"

	// EnvoyMissingAttributeSourceValue is a special attribute value, which can
	// happen when (eg. Envoy's) logger tries to send attribute value, but its
	// source is missing. Eg. In case authz couldn't reach agent, so we know
	// nothing about flowcontrol policies.  Note that this is different case
	// from "just empty", eg. "", "[]" or "{}".
	EnvoyMissingAttributeSourceValue = "-"

	/* HTTP Specific labels. */

	// HTTPStatusCodeLabel describes HTTP status code of the response.
	HTTPStatusCodeLabel = "http.status_code"
	// HTTPRequestContentLength describes length of the HTTP request content in bytes.
	HTTPRequestContentLength = "http.request_content_length"
	// HTTPResponseContentLength describes length of the HTTP response content in bytes.
	HTTPResponseContentLength = "http.response_content_length"
	// HTTPMethodLabel describes HTTP method of the request.
	HTTPMethodLabel = "http.method"
	// HTTPTargetLabel describes HTTP target of the request.
	HTTPTargetLabel = "http.target"
	// HTTPFlavorLabel describes HTTP flavor of the request.
	HTTPFlavorLabel = "http.flavor"
	// HTTPUserAgentLabel describes HTTP user agent of the request.
	HTTPUserAgentLabel = "http.user_agent"
	// HTTPHostLabel describes HTTP host of the request.
	HTTPHostLabel = "http.host"

	/* Labels specific to Envoy. */

	// EnvoyDurationLabel from envoy access logs.
	EnvoyDurationLabel = "DURATION"
	// EnvoyRequestDurationLabel from envoy access logs.
	EnvoyRequestDurationLabel = "REQUEST_DURATION"
	// EnvoyRequestTxDurationLabel from envoy access logs.
	EnvoyRequestTxDurationLabel = "REQUEST_TX_DURATION"
	// EnvoyResponseDurationLabel from envoy access logs.
	EnvoyResponseDurationLabel = "RESPONSE_DURATION"
	// EnvoyResponseTxDurationLabel from envoy access logs.
	EnvoyResponseTxDurationLabel = "RESPONSE_TX_DURATION"
	// EnvoyCallerLabel from envoy access logs.
	EnvoyCallerLabel = "caller"

	/* The following are derived labels that are applied based on contents of check response. */

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
	// WorkloadsLabel describes workloads matched to the traffic.
	WorkloadsLabel = "workloads"
	// DroppingWorkloadsLabel describes workloads dropping the traffic.
	DroppingWorkloadsLabel = "dropping_workloads"
	// FluxMetersLabel describes flux meters matched to the traffic.
	FluxMetersLabel = "flux_meters"
	// FlowLabelKeysLabel describes keys of flow labels matched to the traffic.
	FlowLabelKeysLabel = "flow_label_keys"
	// ClassifiersLabel describes classifiers matched to the traffic.
	ClassifiersLabel = "classifiers"

	/* L3-L4 labels. */

	// HostAddressLabel describes host address of the request.
	HostAddressLabel = "net.host.address"
	// PeerAddressLabel describes peer address of the request.
	PeerAddressLabel = "net.peer.address"
	// HostIPLabel describes host IP address of the request.
	HostIPLabel = "net.host.ip"
	// PeerIPLabel describes peer IP address of the request.
	PeerIPLabel = "net.peer.ip"

	/* SDK specific labels. */

	// FeatureAddressLabel describes feature address of the request.
	FeatureAddressLabel = "feature.ip"
	// FeatureIDLabel describes the ID of the feature.
	FeatureIDLabel = "feature.id"
	// FeatureStatusLabel describes the status of the feature.
	FeatureStatusLabel = "feature.status"

	/* Specific to infra metrics pipeline. */

	// EntityNameLabel describes entity name e.g. pod name.
	EntityNameLabel = "entity_name"

	/* Aperture specific enrichment labels for agent group and service. */

	// AgentGroupLabel describes cluster to which metrics refer.
	AgentGroupLabel = "agent_group"
	// ServicesLabel describes services to which metrics refer.
	ServicesLabel = "services"
)
