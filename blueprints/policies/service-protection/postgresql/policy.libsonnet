local promqlPolicyFn = import '../promql/policy.libsonnet';
local config = import './config.libsonnet';

function(cfg, metadata={}) {
  local params = config + cfg,

  local promqlPolicy = promqlPolicyFn(cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
