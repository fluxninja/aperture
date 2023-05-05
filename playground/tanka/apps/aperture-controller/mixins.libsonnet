local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local extensionEnv = std.extVar('ENABLE_CLOUD_EXTENSION');
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
            endpoint: 'aperture.latest.dev.fluxninja.com' + ':443',
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
            create: extensionEnv,
            value: '2b97802cf7984791919758a537c05ad0',
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
