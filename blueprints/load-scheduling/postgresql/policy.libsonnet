local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg, metadata={}) {
  local policyName = cfg.policy.policy_name,
  local promqlQuery = cfg.policy.promql_query % { policy_name: policyName },

  local updated_cfg = cfg {
    policy+: {
      promql_query: promqlQuery,
    },
  },
  local promqlPolicy = promqlPolicyFn(updated_cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
