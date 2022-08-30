local istioApp = import 'apps/istio/main.libsonnet';

local istioAppMixin =
  istioApp {
    values+: {
      istio+: {},
      envoyfilter+: {
      },
    },
  };

istioAppMixin
