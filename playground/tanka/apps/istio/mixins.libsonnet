local istioApp = import 'apps/istio/main.libsonnet';

local istioAppMixin =
  istioApp {
    values+: {
      base+: {},
      istiod+: {},
      gateway+: {},
      envoyfilter+: {
        envoyFilter+: {
          workloadSelector: {
            labels: {
              'app.kubernetes.io/name': 'demo-app',
            },
          },
        },
      },
    },
  };

istioAppMixin
