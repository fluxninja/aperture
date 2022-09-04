local k = import 'k.libsonnet';

local apertureAgentApp = import 'apps/aperture-agent/main.libsonnet';

local apertureAgentMixin =
  apertureAgentApp {
    values+:: {
      operator+: {
        image: {
          registry: 'docker.io/fluxninja',
          repository: 'aperture-operator',
          tag: 'latest',
        },
      },
      agent+: {
        createUninstallHook: false,
        fluxninjaPlugin+: {
          enabled: false,
        },
        log+: {
          prettyConsole: true,
          nonBlocking: false,
          level: 'debug',
          file: 'default',
        },
        image: {
          registry: '',
          repository: 'docker.io/fluxninja/aperture-agent',
          tag: 'latest',
        },
        sidecar+: {
          enabled: false,
        },
        etcd+: {
          endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
        },
        prometheus+: {
          address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
        },
      },
    },
  };

{
  agent: apertureAgentMixin,
}
