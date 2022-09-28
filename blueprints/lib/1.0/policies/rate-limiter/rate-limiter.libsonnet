local spec = import '../../spec.libsonnet';

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

local policy = spec.v1.Policy;
local component = spec.v1.Component;
local constant = spec.v1.Constant;
local rateLimiter = spec.v1.RateLimiter;
local selector = spec.v1.Selector;
local serviceSelector = spec.v1.ServiceSelector;
local flowSelector = spec.v1.FlowSelector;
local circuit = spec.v1.Circuit;
local override = spec.v1.RateLimiterOverride;
local lazySync = spec.v1.RateLimiterLazySync;
local port = spec.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');


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

  local constants = [
    component.withConstant(constant.new()
                           + constant.withValue($._config.rateLimit)
                           + constant.withOutPorts({ output: rateLimitPort })),
  ],

  local policyDef =
    policy.new()
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval($._config.evaluationInterval)
      + circuit.withComponents(constants + [
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
