local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local pluginEnv = std.extVar('ENABLE_CLOUD_PLUGIN');

local apertureControllerMixin =
  apertureControllerApp {
    values+:: {
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
          fluxninja_plugin+: {
            fluxninja_endpoint: 'aperture.latest.dev.fluxninja.com' + ':443',
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
          plugins+: {
            disabled_plugins: if pluginEnv == 'True' then [] else ['aperture-plugin-fluxninja'],
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
          fluxNinjaPlugin+: {
            create: pluginEnv,
            value: '2b97802cf7984791919758a537c05ad0',
          },
        },
        image: {
          registry: '',
          repository: 'docker.io/fluxninja/aperture-controller',
          tag: 'latest',
        },
      },
    },
  };

{
  controller: apertureControllerMixin,
}
