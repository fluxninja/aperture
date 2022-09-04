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
          repository: 'docker.io/fluxninja/aperture-controller',
          tag: 'latest',
        },
      },
    },
  };

{
  controller: apertureControllerMixin,
}
