local k = import 'k.libsonnet';

local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local apertureControllerMixin =
  apertureControllerApp {
    values+:: {
      operator+: {
        image: {
          registry: 'docker.io/fluxninja',
          repository: 'aperture-operator',
          tag: 'latest',
        },
      },
      controller+: {
        createUninstallHook: false,
        config+: {
          plugins+: {
            disabled_plugins: [
              'aperture-plugin-fluxninja',
            ],
          },
          log+: {
            pretty_console: true,
            non_blocking: false,
            level: 'debug',
            file: 'default',
          },
          etcd+: {
            endpoints: ['http://controller-etcd:2379'],
            lease_ttl: '60s',
          },
          prometheus+: {
            address: 'http://controller-prometheus-server.local:80',
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
