local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';

local policy = spec.v1.Policy;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local rateLimiter = spec.v1.RateLimiter;
local alerter = spec.v1.Alerter;
local port = spec.v1.Port;

function(cfg) {
  local params = cfg,
  local policyDef =
    policy.new()
    + policy.withResources(utils.resources(params.policy.resources).updatedResources)
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval('1s')
      + circuit.withComponents(
        [
          component.withFlowControl(
            flowControl.new()
            + flowControl.withRateLimiter(
              rateLimiter.new()
              + rateLimiter.withInPorts({
                bucket_capacity: port.withConstantSignal(params.policy.rate_limiter.bucket_capacity),
                fill_amount: port.withConstantSignal(params.policy.rate_limiter.fill_amount),
              })
              + rateLimiter.withOutPorts({
                accept_percentage: port.withSignalName('ACCEPT_PERCENTAGE'),
              })
              + rateLimiter.withSelectors(params.policy.rate_limiter.selectors)
              + rateLimiter.withParameters(params.policy.rate_limiter.parameters)
              + rateLimiter.withRequestParameters(params.policy.rate_limiter.request_parameters)
            ),
          ),
          component.withDecider(
            spec.v1.Decider.withOperator('gte')
            + spec.v1.Decider.withInPorts({
              lhs: spec.v1.Port.withSignalName('ACCEPT_PERCENTAGE'),
              rhs: spec.v1.Port.withConstantSignal(90),
            })
            + spec.v1.Decider.withOutPorts({
              output: spec.v1.Port.withSignalName('ACCEPT_PERCENTAGE_ALERT'),
            })
          ),
          component.withAlerter(
            alerter.new()
            + alerter.withInPorts({
              alert: port.withSignalName('ACCEPT_PERCENTAGE_ALERT'),
            })
            + alerter.withParameters(params.policy.rate_limiter.alerter)
          ),
        ] + params.policy.components,
      )
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
