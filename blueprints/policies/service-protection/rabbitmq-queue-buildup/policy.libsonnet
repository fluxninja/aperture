local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config.common + config.policy + cfg,

  local basePolicy = basePolicyFn(cfg).policyDef,

  // Add new components to basePolicy
  local policyDef = basePolicy {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = 'sum(rabbitmq_message_current{rabbitmq_queue_name="%(queue_name)s",state="ready"})' % { queue_name: params.queue_name };
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval('1s')
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.mul(
          spec.v1.Port.withConstantSignal(1),
          spec.v1.Port.withConstantSignal(params.latency_baseliner.queue_buildup_setpoint),
          output=spec.v1.Port.withSignalName('SETPOINT')
        )),
      ],
    },
  },

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
