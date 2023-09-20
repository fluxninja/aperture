local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg, metadata={}) {
  local policyName = cfg.policy.policy_name,
  local promqlQuery = '(sum(postgresql_backends{policy_name="%(policy_name)s",infra_meter_name="postgresql"}}) / sum(postgresql_connection_max{policy_name="%(policy_name)s",infra_meter_name="postgresql"})) * 100' % { policy_name: policyName },

  local updated_cfg = cfg {
    policy+: {
      promql_query: promqlQuery,
      setpoint: 40,
    },
  },
  local promqlPolicy = promqlPolicyFn(updated_cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
