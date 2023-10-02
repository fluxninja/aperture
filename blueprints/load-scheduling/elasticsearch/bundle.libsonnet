local creator = import '../../grafana/dashboard_group.libsonnet';
local utils = import '../common/utils.libsonnet';
local blueprint = import './elasticsearch.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg {
    policy+: {
      promql_query: 'avg(avg_over_time(elasticsearch_node_thread_pool_tasks_queued{thread_pool_name="search"}[30s]))',
      setpoint: c.policy.load_scheduling_core.setpoint,
      overload_condition: 'gt',
    },
  },

  local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {},
  assert !std.objectHas(infraMeters, 'elasticsearch') : 'An infra meter with name elasticsearch already exists. Please choose a different name.',
  local config_with_elasticsearch_infra_meter = updated_cfg {
    policy+: {
      resources+: {
        infra_meters+: {
          elasticsearch: {
            agent_group: if std.objectHas(updated_cfg.policy.elasticsearch, 'agent_group') then updated_cfg.policy.elasticsearch.agent_group else 'default',
            per_agent_group: false,
            receivers: {
              elasticsearch: std.prune(updated_cfg.policy.elasticsearch {
                agent_group: null,
                collection_interval: '10s',
                metrics+: {
                  'elasticsearch.node.operations.current': {
                    enabled: true,
                  },
                  'elasticsearch.node.operations.completed': {
                    enabled: true,
                  },
                  'elasticsearch.node.operations.time': {
                    enabled: true,
                  },
                  'elasticsearch.node.operations.get.completed': {
                    enabled: true,
                  },
                  'elasticsearch.node.operations.get.time': {
                    enabled: true,
                  },
                  'jvm.memory.heap.utilization': {
                    enabled: true,
                  },
                },
              }),
            },
          },
        },
      },
    },
  },

  local p = policy(config_with_elasticsearch_infra_meter),
  local d = creator(p.policyResource, config_with_elasticsearch_infra_meter),

  policies: {
    [std.format('%s-cr.yaml', config_with_elasticsearch_infra_meter.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', config_with_elasticsearch_infra_meter.policy.policy_name)]: p.policyDef,
  },
  dashboards: {
    [std.format('%s.json', config_with_elasticsearch_infra_meter.policy.policy_name)]: d.mainDashboard,
    [std.format('signals-%s.json', config_with_elasticsearch_infra_meter.policy.policy_name)]: d.signalsDashboard,
  } + d.receiverDashboards,
}
