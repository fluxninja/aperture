local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg, metadata={}) {
  local updated_cfg = cfg {
    policy+: {
      promql_query: '(sum(postgresql_backends) / sum(postgresql_connection_max)) * 100',
      setpoint: 40,
    },
  },
  local promqlPolicy = promqlPolicyFn(updated_cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
