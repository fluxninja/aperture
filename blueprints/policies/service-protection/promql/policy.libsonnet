local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config + cfg,

  local basePolicy = basePolicyFn(cfg).policyDef,

  // Add new components to basePolicy
  local policyDef = basePolicy {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = params.policy.promql_query;
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval(params.policy.evaluation_interval)
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withVariable(
          spec.v1.Variable.new()
          + spec.v1.Variable.withDefaultConfig(
            spec.v1.VariableDynamicConfig.new()
            + spec.v1.VariableDynamicConfig.withConstantSignal(
              local s = params.policy.setpoint;
              spec.v1.ConstantSignal.new()
              + spec.v1.ConstantSignal.withValue(s)
            )
          )
          + spec.v1.Variable.withOutPorts({ output: spec.v1.Port.withSignalName('SETPOINT') }),
        ),
      ],
    },
  },

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
