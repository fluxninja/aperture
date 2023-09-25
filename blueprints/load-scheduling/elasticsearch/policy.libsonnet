local promqlPolicyFn = import '../promql/policy.libsonnet';

function(cfg, metadata={}) {
  local updated_cfg = cfg {
    policy+: {
      promql_query: 'avg(avg_over_time(elasticsearch_node_thread_pool_tasks_queued{thread_pool_name="search"}[30s]))',
      setpoint: 250,
      service_protection_core+: {
        overload_condition: 'gt',
      },
    },
  },
  local promqlPolicy = promqlPolicyFn(updated_cfg, metadata),

  policyResource: promqlPolicy.policyResource {
    spec+: promqlPolicy.policyDef,
  },
  policyDef: promqlPolicy.policyDef,
}
