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
  local postgresql = std.prune(c.policy.postgresql { agent_group: null }),

  local updatedConfig =
    c {
      policy+: {
        resources+: {
          telemetry_collectors+: [
            {
              agent_group: agent_group,
              infra_meters: config.kube_stat_infra_meter {
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
    if c.policy.service_protection_core.cpu_overload_confirmation_threshold != '' || c.policy.service_protection_core.cpu_overload_confirmation_threshold != null then
      local updatedTelemetryCollectors = std.map(
        function(c) if c.agent_group == agent_group then
          c { infra_meters+: config.kube_stat_infra_meter }
        else c,
        c.policy.resources.telemetry_collectors,
      );
      {
        policy+: {
          resources+: {
            telemetry_collectors: updatedTelemetryCollectors,
          },
          service_protection_core+: {
            overload_confirmations+: [
              {
                query_string: 'max(k8s_pod_cpu_utilization_ratio{k8s_statefulset_name="hasura-postgresql"})',
                threshold: c.policy.service_protection_core.cpu_overload_confirmation_threshold,
                operator: 'gte',
              },
            ],
          },
        },
      }
    else {},

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
