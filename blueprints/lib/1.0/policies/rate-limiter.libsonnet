local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

// A set of defaults used by the policy that can be overridden when instantiating it
local defaults = {
  policyName: error 'policyName must be set',
  evaluationInterval: '0.5s',
  serviceSelector: {
    agentGroup: 'default',
    service: error 'policy serviceSelector.service is required',
    controlPoint: {
      traffic: 'ingress',
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
local constant = aperture.v1.Constant;
local rateLimiter = aperture.v1.RateLimiter;
local selector = aperture.v1.Selector;
local circuit = aperture.v1.Circuit;
local override = aperture.v1.RateLimiterOverride;
local lazySync = aperture.v1.RateLimiterLazySync;
local port = aperture.v1.Port;

local rateLimitPort = port.new() + port.withSignalName('RATE_LIMIT');


function(params) {
  _config:: defaults + params,

  local svcSelector =
    selector.new()
    + selector.withAgentGroup($._config.serviceSelector.agentGroup)
    + selector.withService($._config.serviceSelector.service)
    + selector.withControlPoint({ traffic: $._config.serviceSelector.controlPoint.traffic }),

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
