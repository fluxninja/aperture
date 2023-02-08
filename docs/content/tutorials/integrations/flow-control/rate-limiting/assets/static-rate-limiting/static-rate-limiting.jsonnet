local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';


local StaticRateLimiting = aperture.policies.StaticRateLimiting.policy;
local policy = aperture.spec.v1.Policy;
local component = aperture.spec.v1.Component;
local flowControl = aperture.spec.v1.FlowControl;
local rateLimiter = aperture.spec.v1.RateLimiter;
local rateLimiterParameters = aperture.spec.v1.RateLimiterParameters;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local circuit = aperture.spec.v1.Circuit;
local resources = aperture.spec.v1.Resources;
local port = aperture.spec.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint('ingress')
  );

local policyResource = StaticRateLimiting({
  policy_name: 'static-rate-limiting',
  evaluation_interval: '300s',
  rate_limiter+: {
    flow_selector: svcSelector,
    rate_limit: 120.0,
    parameters+: {
      label_key: 'http.request.header.user_id',
      limit_reset_interval: '60s',
    },
  },

}).policyResource;

policyResource
