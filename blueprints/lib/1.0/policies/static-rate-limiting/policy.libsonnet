local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local resources = spec.v1.Resources;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local rateLimiter = spec.v1.RateLimiter;
local override = spec.v1.RateLimiterOverride;
local lazySync = spec.v1.RateLimiterLazySync;
local port = spec.v1.Port;
local constantSignal = spec.v1.ConstantSignal;

function(params) {
  _config:: config.common + config.policy + params,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval($._config.evaluation_interval)
      + circuit.withComponents([
        component.withFlowControl(
          flowControl.new()
          + flowControl.withRateLimiter(
            rateLimiter.new()
            + rateLimiter.withInPorts({ limit: port.withConstantSignal($._config.rate_limiter.rate_limit) })
            + rateLimiter.withFlowSelector($._config.rate_limiter.flow_selector)
            + rateLimiter.withParameters($._config.rate_limiter.parameters)
            + rateLimiter.withDynamicConfigKey('rate_limiter')
          ),
        ),
      ]),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: $._config.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
    dynamicConfig: {
      rate_limiter: $._config.rate_limiter.dynamic_config,
    },
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
