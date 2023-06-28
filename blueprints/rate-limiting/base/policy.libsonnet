local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local rateLimiter = spec.v1.RateLimiter;
local port = spec.v1.Port;

function(cfg, metadata={}) {
  local params = config + cfg,
  local policyDef =
    policy.new()
    + policy.withResources(utils.resources(params.policy.resources).updatedResources)
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval('1s')
      + circuit.withComponents([
        component.withFlowControl(
          flowControl.new()
          + flowControl.withRateLimiter(
            rateLimiter.new()
            + rateLimiter.withInPorts({
              bucket_capacity: port.withConstantSignal(params.policy.rate_limiter.bucket_capacity),
              fill_amount: port.withConstantSignal(params.policy.rate_limiter.fill_amount),
            })
            + rateLimiter.withSelectors(params.policy.rate_limiter.selectors)
            + rateLimiter.withParameters(params.policy.rate_limiter.parameters)
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
      annotations: {
        [if std.objectHas(metadata, 'values') then 'fluxninja.com/values']: metadata.values,
        [if std.objectHas(metadata, 'blueprints_uri') then 'fluxninja.com/blueprints-uri']: metadata.blueprints_uri,
        [if std.objectHas(metadata, 'blueprint_name') then 'fluxninja.com/blueprint-name']: metadata.blueprint_name,
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
