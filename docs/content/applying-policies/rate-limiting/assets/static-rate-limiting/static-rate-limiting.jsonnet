local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';


local RateLimiting = aperture.policies.RateLimiting.policy;
local selector = aperture.spec.v1.Selector;

local svcSelectors = [
  selector.new()
  + selector.withControlPoint('ingress')
  + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
  + selector.withAgentGroup('default'),
];

local policyResource = RateLimiting({
  policy+: {
    policy_name: 'static-rate-limiting',
    rate_limiter+: {
      selectors: svcSelectors,
      bucket_capacity: 40,
      fill_amount: 2,
      parameters+: {
        label_key: 'http.request.header.user_id',
        interval: '1s',
      },
    },
  },
}).policyResource;

policyResource
