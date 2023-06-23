local spec = import '../../../spec.libsonnet';
local basePolicyFn = import '../base/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, metadata={}) {
  local params = config + cfg,

  local basePolicy = basePolicyFn(cfg, metadata),

  // Add new components to basePolicy
  local policyDef = basePolicy.policyDef {
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
          + spec.v1.Variable.withConstantOutput(
            local s = params.policy.setpoint;
            spec.v1.ConstantSignal.new()
            + spec.v1.ConstantSignal.withValue(s)
          )
          + spec.v1.Variable.withOutPorts({ output: spec.v1.Port.withSignalName('SETPOINT') }),
        ),
      ],
    },
  },

  policyResource: basePolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
