local aperture = import '../../lib/1.0/main.libsonnet';
local bundle = aperture.blueprints.RateLimiter.bundle;

local Override = aperture.spec.v1.RateLimiterOverride;
local LazySync = aperture.spec.v1.RateLimiterLazySync;

local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;


local svcSelector = flowSelector.new()
                    + flowSelector.withServiceSelector(
                      serviceSelector.new()
                      + serviceSelector.withAgentGroup('default')
                      + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                    )
                    + flowSelector.withFlowMatcher(
                      flowMatcher.new()
                      + flowMatcher.withControlPoint('ingress')
                    );

local config = {
  common+: {
    policyName: 'example',
  },
  policy+: {
    rateLimiterFlowSelector: svcSelector,
    rateLimit: '50.0',
    labelKey: 'http.request.header.user_type',
    limitResetInterval: '1s',
    overrides: [
      Override.new()
      + Override.withLabelValue('gold')
      + Override.withLimitScaleFactor(1),
    ],
  },
};

bundle { _config+:: config }
