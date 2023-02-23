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

function(cfg) {
  local params = config.common + config.policy + cfg,
  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withClassifiers(params.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval('1s')
      + circuit.withComponents([
        component.withFlowControl(
          flowControl.new()
          + flowControl.withRateLimiter(
            rateLimiter.new()
            + rateLimiter.withInPorts({ limit: port.withConstantSignal(params.rate_limiter.rate_limit) })
            + rateLimiter.withFlowSelector(params.rate_limiter.flow_selector)
            + rateLimiter.withParameters(params.rate_limiter.parameters)
            + rateLimiter.withDynamicConfigKey('rate_limiter')
            + rateLimiter.withDefaultConfig(params.rate_limiter.default_config)
          ),
        ),
      ]),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
