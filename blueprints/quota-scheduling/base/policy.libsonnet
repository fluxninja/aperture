local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';

local policy = spec.v1.Policy;
local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local quotaScheduler = spec.v1.QuotaScheduler;
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
            + flowControl.withQuotaScheduler(
              quotaScheduler.new()
              + quotaScheduler.withInPorts({
                bucket_capacity: port.withConstantSignal(params.policy.quota_scheduler.bucket_capacity),
                fill_amount: port.withConstantSignal(params.policy.quota_scheduler.fill_amount),
              })
              + quotaScheduler.withOutPorts({
                accept_percentage: port.withSignalName('ACCEPT_PERCENTAGE'),
              })
              + quotaScheduler.withSelectors(params.policy.quota_scheduler.selectors)
              + quotaScheduler.withRateLimiter(params.policy.quota_scheduler.rate_limiter)
              + quotaScheduler.withScheduler(params.policy.quota_scheduler.scheduler)
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
            + alerter.withParameters(params.policy.quota_scheduler.alerter)
          ),
        ] + params.policy.components,
      ),
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
