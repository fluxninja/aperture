local spec = import '../../spec.libsonnet';
local commonPolicyFn = import '../common/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, params={}, metadata={}) {
  local updatedConfig = config + cfg,

  local commonPolicy = commonPolicyFn(cfg, params, metadata),

  // Add new components to commonPolicy
  local policyDef = commonPolicy.policyDef {
    circuit+: {
      components+: [
        spec.v1.Component.withQuery(
          spec.v1.Query.new()
          + spec.v1.Query.withPromql(
            local q = updatedConfig.policy.promql_query;
            spec.v1.PromQL.new()
            + spec.v1.PromQL.withQueryString(q)
            + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=updatedConfig.policy.evaluation_interval)
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
        spec.v1.Component.withVariable(
          spec.v1.Variable.new()
          + spec.v1.Variable.withConstantOutput(
            local s = updatedConfig.policy.setpoint;
            spec.v1.ConstantSignal.new()
            + spec.v1.ConstantSignal.withValue(s)
          )
          + spec.v1.Variable.withOutPorts({ output: spec.v1.Port.withSignalName('SETPOINT') }),
        ),
      ],
    },
  },

  policyResource: commonPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
