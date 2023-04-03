local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';


local StaticRateLimiting = aperture.policies.StaticRateLimiting.policy;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint('ingress')
  );

local policyResource = StaticRateLimiting({
  policy_name: 'static-rate-limiting',
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
