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
    c {
      policy+: {
        resources+: {
          telemetry_collectors+: [
            {
              agent_group: random_string,
              infra_meters+: {
                postgresql: {
                  receivers: {
                    postgresql: postgresql,
                  },
                },
              },
            },
          ],
        },
      },
    },

  local finalConfig =
    updatedConfig +
    if std.objectHas(params.policy, 'service_protection_core')
       && std.objectHas(params.policy.service_protection_core, 'cpu_overload_confirmation') then
      assert updatedConfig.policy.service_protection_core.cpu_overload_confirmation.query_string != ''
             && updatedConfig.policy.service_protection_core.cpu_overload_confirmation.query_string != null : 'query_string must be set for policy.service_protection_core.cpu_overload_confirmation';
      assert updatedConfig.policy.service_protection_core.cpu_overload_confirmation.threshold != ''
             && updatedConfig.policy.service_protection_core.cpu_overload_confirmation.threshold != null : 'threshold must be set for policy.service_protection_core.cpu_overload_confirmation';
      assert updatedConfig.policy.service_protection_core.cpu_overload_confirmation.operator != ''
             && updatedConfig.policy.service_protection_core.cpu_overload_confirmation.operator != null : 'operator must be set for policy.service_protection_core.cpu_overload_confirmation';
      local updatedTelemetryCollectors = std.map(
        function(collector) if collector.agent_group == random_string then
          collector { infra_meters+: config.kubeletstats_infra_meter, agent_group: agent_group }
        else collector,
        updatedConfig.policy.resources.telemetry_collectors,
      );
      {
        policy+: {
          resources+: {
            telemetry_collectors: updatedTelemetryCollectors,
          },
          service_protection_core+: {
            overload_confirmations+: [
              updatedConfig.policy.service_protection_core.cpu_overload_confirmation,
            ],
          },
        },
      }
    else {},

  local metadataWrapper = metadata { values: std.toString(params) },
  local p = policy(finalConfig, metadataWrapper),
  local d = dashboard(finalConfig),

  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
  dashboards: {
    [std.format('%s.json', c.policy.policy_name)]: d.dashboard,
  },
}
