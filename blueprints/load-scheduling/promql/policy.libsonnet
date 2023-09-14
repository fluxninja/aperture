local spec = import '../../spec.libsonnet';
local commonPolicyFn = import '../common-range/policy.libsonnet';
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
            + spec.v1.PromQL.withEvaluationInterval(evaluation_interval='10s')
            + spec.v1.PromQL.withOutPorts({ output: spec.v1.Port.withSignalName('SIGNAL') }),
          ),
        ),
      ],
    },
  },

  policyResource: commonPolicy.policyResource {
    spec+: policyDef,
  },
  policyDef: policyDef,
}
