local k = import 'k.libsonnet';

local apertureOperatorApp = import 'apps/aperture-operator/main.libsonnet';

local apertureOperatorMixin =
  apertureOperatorApp {
    values+:: {
      global+: {
        istioNamespace: 'aperture-system',
      },
      operator+: {
        image: {
          registry: 'gcr.io/devel-309501/cf-fn',
          repository: 'aperture-operator',
          tag: 'latest',
        },
      },
      aperture+: {
        createUninstallHook: false,
        fluxninjaPlugin+: {
          enabled: false,
        },
        agent+: {
          log+: {
            prettyConsole: true,
            nonBlocking: false,
            level: 'debug',
            file: 'default',
          },
          image: {
            registry: '',
            repository: 'gcr.io/devel-309501/cf-fn/aperture-agent',
            tag: 'latest',
          },
        },
        controller+: {
          log+: {
            prettyConsole: true,
            nonBlocking: false,
            level: 'debug',
            file: 'default',
          },
          image: {
            registry: '',
            repository: 'gcr.io/devel-309501/cf-fn/aperture-controller',
            tag: 'latest',
          },
        },
        sidecar+: {
          enabled: false,
        },
      },
    },
  };

apertureOperatorMixin
