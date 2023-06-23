local blueprint = import './postgresql.libsonnet';

local policy = blueprint.policy;
local dashboard = blueprint.dashboard;
local config = blueprint.config;

function(params, metadata={}) {
  // make sure param object contains fields that are in config
  local extra_keys = std.setDiff(std.objectFields(params), std.objectFields(config)),
  assert std.length(extra_keys) == 0 : 'Unknown keys in params: ' + extra_keys,

  local c = std.mergePatch(config, params),
  local agent_group = c.policy.postgresql.agent_group,
  local random_string = std.substr(std.md5(std.toString(params)), 0, 5),
  local postgresql = std.prune(c.policy.postgresql { agent_group: null }),

  local updatedConfig =
    local addCPUOverloadConfirmation =
      if std.objectHas(params.policy, 'service_protection_core')
         && std.objectHas(params.policy.service_protection_core, 'cpu_overload_confirmation') then
        assert params.policy.service_protection_core.cpu_overload_confirmation.query_string != ''
               && params.policy.service_protection_core.cpu_overload_confirmation.query_string != null : 'query_string must be set for policy.service_protection_core.cpu_overload_confirmation';
        assert params.policy.service_protection_core.cpu_overload_confirmation.threshold != null : 'threshold must be set for policy.service_protection_core.cpu_overload_confirmation';
        assert params.policy.service_protection_core.cpu_overload_confirmation.operator != ''
               && params.policy.service_protection_core.cpu_overload_confirmation.operator != null : 'operator must be set for policy.service_protection_core.cpu_overload_confirmation';
        true else false;
    c {
      policy+: {
        resources+: {
          infra_meters+: {
            postgresql: {
              agent_group: agent_group,
              receivers: {
                postgresql: postgresql,
              },
            },
          } + if addCPUOverloadConfirmation then config.kubeletstats_infra_meter(agent_group) else {},
        },
        service_protection_core+: if addCPUOverloadConfirmation then {
          overload_confirmations+: [
            updatedConfig.policy.service_protection_core.cpu_overload_confirmation,
          ],
        } else {},
      },
    },

  local metadataWrapper = metadata { values: std.toString(params) },
  local p = policy(updatedConfig, metadataWrapper),
  local d = dashboard(updatedConfig),

  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
  dashboards: {
    [std.format('%s.json', c.policy.policy_name)]: d.dashboard,
  },
}
