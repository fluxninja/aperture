local creator = import '../../../grafana/creator.libsonnet';
local blueprint = import './postgresql.libsonnet';

local policy = blueprint.policy;
local config = blueprint.config;

local kubeletstats_infra_meter = function(agent_group='default') {
  kubeletstats: {
    agent_group: agent_group,
    pipeline: {
      processors: [
        'k8sattributes',
      ],
      receivers: [
        'kubeletstats',
      ],
    },
    processors: {
      k8sattributes: {
        auth_type: 'serviceAccount',
        passthrough: false,
        extract: {
          metadata: [
            'k8s.cronjob.name',
            'k8s.daemonset.name',
            'k8s.deployment.name',
            'k8s.job.name',
            'k8s.namespace.name',
            'k8s.node.name',
            'k8s.pod.name',
            'k8s.pod.uid',
            'k8s.replicaset.name',
            'k8s.statefulset.name',
            'k8s.container.name',
          ],
        },
        pod_association: [
          {
            sources: [
              {
                from: 'resource_attribute',
                name: 'k8s.pod.uid',
              },
            ],
          },
        ],
      },
    },
    receivers: {
      kubeletstats: {
        collection_interval: '15s',
        auth_type: 'serviceAccount',
        endpoint: 'https://${NODE_NAME}:10250',
        insecure_skip_verify: true,
        metric_groups: [
          'pod',
          'container',
        ],
      },
    },
  },
};

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
              per_agent_group: true,
              receivers: {
                postgresql: postgresql,
              },
            },
          } + if addCPUOverloadConfirmation then kubeletstats_infra_meter(agent_group) else {},
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
  local d = creator(p.policyResource, updatedConfig),

  policies: {
    [std.format('%s-cr.yaml', c.policy.policy_name)]: p.policyResource,
    [std.format('%s.yaml', c.policy.policy_name)]: p.policyDef { metadata: metadataWrapper },
  },
  dashboards: {
    [std.format('%s.json', c.policy.policy_name)]: d.dashboard,
  },
}
