local istioApp = import 'apps/istio/main.libsonnet';

local istioAppMixin =
  istioApp {
    values+: {
      base+: {},
      istiod+: {},
      gateway+: {},
      envoyfilter+: {},
    },
  };

istioAppMixin
