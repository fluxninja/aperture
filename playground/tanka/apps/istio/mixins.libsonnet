local istioApp = import 'apps/istio/main.libsonnet';

local istioAppMixin =
  istioApp {
    values+: {
      base+: {},
      istiod+: {},
      gateway+: {},
      envoyfilter+: {
        envoyFilter+: {
          authzPort: 80,
        },
      },
    },
  };

istioAppMixin
