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
            non_blocking: true,
            level: 'debug',
          },
          etcd+: {
            endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
          },
          prometheus+: {
            address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
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
