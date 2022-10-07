local spec = import '../../spec.libsonnet';

// A set of defaults used by the policy that can be overridden when instantiating it
local defaults = {
  policyName: error 'policyName must be set',
  evaluationInterval: '0.5s',
  rateLimiterSelector: error 'rateLimiterSelector must be set',
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
local rateLimiter = spec.v1.RateLimiter;
local circuit = spec.v1.Circuit;
local dynamicConfig = spec.v1.RateLimiterDynamicConfig;
local override = spec.v1.RateLimiterOverride;
local lazySync = spec.v1.RateLimiterLazySync;
local port = spec.v1.Port;

function(params) {
  _config:: defaults + params,

  local policyDef =
    policy.new()
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval($._config.evaluationInterval)
      + circuit.withComponents([
        component.withRateLimiter(
          rateLimiter.new()
          + rateLimiter.withInPorts({ limit: port.withConstantValue($._config.rateLimit) })
          + rateLimiter.withSelector($._config.rateLimiterSelector)
          + rateLimiter.withLimitResetInterval($._config.limitResetInterval)
          + rateLimiter.withLabelKey($._config.labelKey)
          + rateLimiter.withInitConfig(
            dynamicConfig.new()
            + dynamicConfig.withOverrides($._config.overrides)
          )
          + rateLimiter.withLazySync(lazySync.new()
                                     + lazySync.withEnabled($._config.lazySync.enabled)
                                     + lazySync.withNumSync($._config.lazySync.numSync))
        ),
      ]),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: $._config.policyName,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
}
