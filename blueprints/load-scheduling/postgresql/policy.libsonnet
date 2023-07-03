local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg, metadata={}) {
  local promqlPolicy = promqlPolicyFn(cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
