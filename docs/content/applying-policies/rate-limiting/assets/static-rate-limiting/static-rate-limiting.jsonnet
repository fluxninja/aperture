local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';


local StaticRateLimiting = aperture.policies.StaticRateLimiting.policy;
local selector = aperture.spec.v1.Selector;

local svcSelectors = [
  selector.new()
  + selector.withControlPoint('ingress')
  + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
  + selector.withAgentGroup('default'),
];

local policyResource = StaticRateLimiting({
  policy+: {
    policy_name: 'static-rate-limiting',
    rate_limiter+: {
      selectors: svcSelectors,
      rate_limit: 120.0,
      parameters+: {
        label_key: 'http.request.header.user_id',
        limit_reset_interval: '60s',
      },
    },
  },
}).policyResource;

policyResource
