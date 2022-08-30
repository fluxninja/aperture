local k = import 'k.libsonnet';

local apertureOperatorApp = import 'apps/aperture-operator/main.libsonnet';

local apertureOperatorMixin =
  apertureOperatorApp {
    values+:: {
      operator+: {
        image: {
          registry: 'gcr.io/devel-309501/cf-fn',
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
          repository: 'gcr.io/devel-309501/cf-fn/aperture-agent',
          tag: 'latest',
        },
        sidecar+: {
          enabled: false,
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
          repository: 'gcr.io/devel-309501/cf-fn/aperture-controller',
          tag: 'latest',
        },
      },
    },
  };

apertureOperatorMixin
