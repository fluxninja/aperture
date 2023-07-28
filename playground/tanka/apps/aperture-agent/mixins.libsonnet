local apertureAgentApp = import 'apps/aperture-agent/main.libsonnet';

local extensionEnvStr = std.extVar('CLOUD_EXTENSION');
local extensionEnv = if extensionEnvStr != '' then std.parseYaml(extensionEnvStr) else {};
local cloudController = if std.objectHas(extensionEnv, 'cloud_controller') then extensionEnv.cloud_controller else false;
local apiKey = if std.objectHas(extensionEnv, 'api_key') then extensionEnv.api_key else '';
local endpoint = if std.objectHas(extensionEnv, 'endpoint') then extensionEnv.endpoint else '';
local enabled = if std.objectHas(extensionEnv, 'enabled') then extensionEnv.enabled else false;
local agentGroup = if std.objectHas(extensionEnv, 'agent_group') then extensionEnv.agent_group else '';
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local agentValues = if std.objectHas(values, 'agent') then values.agent else {};

local apertureAgentMixin =
  apertureAgentApp {
    values+:: std.mergePatch({
      operator+: {
        image: {
          registry: 'docker.io/fluxninja',
          repository: 'aperture-operator',
          tag: 'latest',
        },
      },
      agent+: {
        createUninstallHook: false,
        config+: {
          agent_info+: {
            agent_group: if cloudController then agentGroup else 'default',
          },
          fluxninja+: {
            enable_cloud_controller: cloudController,
            endpoint: if enabled || cloudController then endpoint else '',
            client+: {
              grpc+: {
                insecure: false,
                tls+: {
                  insecure_skip_verify: true,
                },
              },
              http+: {
                tls+: {
                  insecure_skip_verify: true,
                },
              },
            },
          },
          log+: {
            pretty_console: true,
            non_blocking: true,
            level: 'info',
          },
          etcd+: if !cloudController then {
            endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
          } else {},
          prometheus+: if !cloudController then {
            address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
          } else {},
          flow_control+: {
            preview_service+: {
              enabled: true,
            },
          },
          agent_functions+: if !cloudController then {
            endpoints: ['aperture-controller.aperture-controller.svc.cluster.local:8080'],
          } else {},
        },
        secrets+: {
          fluxNinjaExtension+: {
            create: enabled || cloudController,
            value: if enabled || cloudController then apiKey else '',
          },
        },
        image: {
          registry: '',
          repository: 'docker.io/fluxninja/aperture-agent',
          tag: 'latest',
        },
        sidecar+: {
          enabled: false,
        },
      },
    }, agentValues),
  };

{
  agent: apertureAgentMixin,
}
