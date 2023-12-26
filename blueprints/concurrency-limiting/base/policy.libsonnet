local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';

local policy = spec.v1.Policy;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local concurrencyLimiter = spec.v1.ConcurrencyLimiter;
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
            + flowControl.withConcurrencyLimiter(
              concurrencyLimiter.new()
              + concurrencyLimiter.withInPorts({
                max_concurrency: port.withConstantSignal(params.policy.concurrency_limiter.max_concurrency),
              })
              + concurrencyLimiter.withOutPorts({
                accept_percentage: port.withSignalName('ACCEPT_PERCENTAGE'),
              })
              + concurrencyLimiter.withSelectors(params.policy.concurrency_limiter.selectors)
              + concurrencyLimiter.withParameters(params.policy.concurrency_limiter.parameters)
              + concurrencyLimiter.withRequestParameters(params.policy.concurrency_limiter.request_parameters)
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
              signal: port.withSignalName('ACCEPT_PERCENTAGE_ALERT'),
            })
            + alerter.withParameters(params.policy.concurrency_limiter.alerter)
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
