local istioApp = import 'apps/istio/main.libsonnet';

local istioAppMixin =
  istioApp {
    values+: {
      base+: {},
      istiod+: {},
      gateway+: {},
      envoyfilter+: {
        apertureAgent+: {
          authzPort: 8080,
        },
      },
    },
  };

istioAppMixin
