local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg) {
  local promqlPolicy = promqlPolicyFn(cfg),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
