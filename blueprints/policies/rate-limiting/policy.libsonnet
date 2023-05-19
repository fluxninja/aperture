local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local resources = spec.v1.Resources;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local flowControlResources = spec.v1.FlowControlResources;
local rateLimiter = spec.v1.RateLimiter;
local port = spec.v1.Port;

function(cfg) {
  local params = config + cfg,
  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFlowControl(flowControlResources.new()
                                                       + flowControlResources.withClassifiers(params.policy.classifiers)))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval('1s')
      + circuit.withComponents([
        component.withFlowControl(
          flowControl.new()
          + flowControl.withRateLimiter(
            rateLimiter.new()
            + rateLimiter.withInPorts({ limit: port.withConstantSignal(params.policy.rate_limiter.rate_limit) })
            + rateLimiter.withSelectors(params.policy.rate_limiter.selectors)
            + rateLimiter.withParameters(params.policy.rate_limiter.parameters)
            + rateLimiter.withDynamicConfigKey('rate_limiter')
            + rateLimiter.withDefaultConfig(params.policy.rate_limiter.default_config)
          ),
        ),
      ]),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
