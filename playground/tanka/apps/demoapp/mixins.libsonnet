local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';

local latencyGradientPolicy = import '../../../../blueprints/lib/1.0/policies/latency-gradient.libsonnet';
local aperture = import '../../../../blueprints/libsonnet/1.0/main.libsonnet';
local demoApp = import 'apps/demoapp/main.libsonnet';

local demoappMixin =
  demoApp {
    values+: {
      replicaCount: 2,
      simplesrv+: {
        image: {
          repository: 'docker.io/fluxninja/demo-app',
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

{
  demoapp: demoappMixin,
}
