local k = import 'k.libsonnet';

local demoApp = import 'apps/demoapp/main.libsonnet';

local demoappMixin =
  demoApp {
    values+: {
      simplesrv+: {
        image: {
          repository: 'gcr.io/devel-309501/cf-fn/demo-app',
          tag: 'test',
        },
      },
      resources+: {
        limits+: {
          cpu: '100m',
          memory: '128Mi',
        },
        requests+: {
          cpu: '100m',
          memory: '128Mi',
        },
      },
    },
  };

demoappMixin
