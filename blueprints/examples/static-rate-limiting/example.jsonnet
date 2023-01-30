local aperture = import '../../lib/1.0/main.libsonnet';
local bundle = aperture.policies.StaticRateLimiting.bundle;

local override = aperture.spec.v1.RateLimiterOverride;

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
    policy_name: 'example',
  },
  policy+: {
    rate_limiter+: {
      flow_selector: svcSelector,
      rate_limit: '50.0',
      parameters+: {
        label_key: 'http.request.header.user_type',
        limit_reset_interval: '1s',
      },
      dynamic_config: {
        overrides: [
          override.new()
          + override.withLabelValue('gold')
          + override.withLimitScaleFactor(1),
        ],
      },
    },
  },
};

bundle { _config+:: config }
