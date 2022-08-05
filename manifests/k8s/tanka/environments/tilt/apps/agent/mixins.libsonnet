local k = import 'k.libsonnet';

local agentApp = import 'apps/agent/main.libsonnet';

local agentMixin =
  agentApp {
    values+:: {
      global+: {
        istioNamespace: 'aperture-system',
      },
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
          registry: 'gcr.io/devel-309501/cf-fn',
          repository: 'aperture-agent',
          tag: 'latest',
        },
      },
      agentController+: {
        log+: {
          prettyConsole: true,
          nonBlocking: false,
          level: 'debug',
          file: 'default',
        },
        image: {
          registry: 'gcr.io/devel-309501/cf-fn',
          repository: 'aperture-controller',
          tag: 'latest',
        },
      },
    },
  };

agentMixin
