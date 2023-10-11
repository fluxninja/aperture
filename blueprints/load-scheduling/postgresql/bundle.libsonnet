local utils = import '../common/utils.libsonnet';
local blueprint = import './postgresql.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

function(params) {
  local c = std.mergePatch(config, params),

  local policyName = c.policy.policy_name,
  local promqlQuery = '(sum(postgresql_backends{policy_name="%(policy_name)s",infra_meter_name="postgresql"}) / sum(postgresql_connection_max{policy_name="%(policy_name)s",infra_meter_name="postgresql"})) * 100' % { policy_name: policyName },

  local updated_cfg = utils.add_kubelet_overload_confirmations(c).updated_cfg {
    policy+: {
      promql_query: promqlQuery,
      setpoint: c.policy.load_scheduling_core.setpoint,
      overload_condition: 'gt',
    },
  },

  local infraMeters = if std.objectHas(c.policy.resources, 'infra_meters') then c.policy.resources.infra_meters else {},
  assert !std.objectHas(infraMeters, 'postgresql') : 'An infra meter with name postgresql already exists. Please choose a different name.',
  local config_with_postgresql_infra_meter = updated_cfg {
    policy+: {
      resources+: {
        infra_meters+: {
          postgresql: {
            agent_group: if std.objectHas(updated_cfg.policy.postgresql, 'agent_group') then updated_cfg.policy.postgresql.agent_group else 'default',
            per_agent_group: true,
            receivers: {
              postgresql: std.prune(updated_cfg.policy.postgresql { agent_group: null, collection_interval: '10s' }),
            },
          },
        },
      },
    },
  },

  local p = policy(config_with_postgresql_infra_meter),
  policies: {
    [std.format('%s-cr.yaml', config_with_postgresql_infra_meter.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', config_with_postgresql_infra_meter.policy.policy_name)]: p.policyDef,
  },
}
