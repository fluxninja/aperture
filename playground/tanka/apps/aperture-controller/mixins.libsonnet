local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local extensionEnvStr = std.extVar('CLOUD_EXTENSION');
local extensionEnv = if extensionEnvStr != '' then std.parseYaml(extensionEnvStr) else {};
local apiKey = if std.objectHas(extensionEnv, 'api_key') then extensionEnv.api_key else '';
local endpoint = if std.objectHas(extensionEnv, 'endpoint') then extensionEnv.endpoint else '';
local enabled = if std.objectHas(extensionEnv, 'enabled') then extensionEnv.enabled else false;
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local controllerValues = if std.objectHas(values, 'controller') then values.controller else {};

local apertureControllerMixin =
  apertureControllerApp {
    values+:: std.mergePatch({
      operator+: {
        image: {
          registry: 'docker.io/fluxninja',
          repository: 'aperture-operator',
          tag: 'latest',
          pullPolicy: 'IfNotPresent',
        },
      },
      controller+: {
        createUninstallHook: false,

        config+: {
          fluxninja+: {
            endpoint: if enabled then endpoint else '',
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
          etcd+: {
            endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
          },
          prometheus+: {
            address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
          },
        },
        secrets+: {
          fluxNinjaExtension+: {
            create: enabled,
            value: if enabled then apiKey else '',
          },
        },
        image: {
          registry: '',
          repository: 'docker.io/fluxninja/aperture-controller',
          tag: 'latest',
        },
      },
    }, controllerValues),
  };

{
  controller: apertureControllerMixin,
}
