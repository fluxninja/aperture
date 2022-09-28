local aperture = import '../../../libsonnet/1.0/main.libsonnet';

// A set of defaults used by the policy that can be overridden when instantiating it
local defaults = {
  policyName: error 'policyName must be set',
  evaluationInterval: '0.5s',
  selector: {
    serviceSelector: {
      agentGroup: 'default',
      service: error 'policy serviceSelector.service is required',
    },
    flowSelector: {
      controlPoint: {
        traffic: 'ingress',
      },
    },
  },
  rateLimit: '50.0',
  labelKey: error 'policy labelKey is required',
  limitResetInterval: '1s',
  overrides: [],
  lazySync: {
    enabled: true,
    numSync: 10,
  },
};

local policy = aperture.v1.Policy;
local component = aperture.v1.Component;
local rateLimiter = aperture.v1.RateLimiter;
local selector = aperture.v1.Selector;
local serviceSelector = aperture.v1.ServiceSelector;
local flowSelector = aperture.v1.FlowSelector;
local circuit = aperture.v1.Circuit;
local override = aperture.v1.RateLimiterOverride;
local lazySync = aperture.v1.RateLimiterLazySync;
local inPort = aperture.v1.InPort;

local rateLimitPort = inPort.new() + inPort.withConstantValue(defaults.rateLimit);

function(params) {
  _config:: defaults + params,

  local svcSelector =
    selector.new()
    + selector.withServiceSelector(
      serviceSelector.new()
      + serviceSelector.withAgentGroup($._config.selector.serviceSelector.agentGroup)
      + serviceSelector.withService($._config.selector.serviceSelector.service)
    )
    + selector.withFlowSelector(
      flowSelector.new()
      + flowSelector.withControlPoint({ traffic: $._config.selector.flowSelector.controlPoint.traffic })
    ),

  local policyDef =
    policy.new()
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval($._config.evaluationInterval)
      + circuit.withComponents([
        component.withRateLimiter(
          rateLimiter.new()
          + rateLimiter.withInPorts({ limit: rateLimitPort })
          + rateLimiter.withSelector(svcSelector)
          + rateLimiter.withLimitResetInterval($._config.limitResetInterval)
          + rateLimiter.withLabelKey($._config.labelKey)
          + rateLimiter.withOverrides($._config.overrides)
          + rateLimiter.withLazySync(lazySync.new()
                                     + lazySync.withEnabled($._config.lazySync.enabled)
                                     + lazySync.withNumSync($._config.lazySync.numSync))
        ),
      ]),
    ),
  policy: policyDef,
}
