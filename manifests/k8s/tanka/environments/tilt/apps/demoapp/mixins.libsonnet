local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local k = import 'github.com/jsonnet-libs/k8s-libsonnet/1.22/main.libsonnet';

local demoApp = import 'apps/demoapp/main.libsonnet';
local latencyGradientPolicy = import 'github.com/fluxninja/aperture-blueprints/lib/1.0/policies/latency-gradient.libsonnet';

local demoappMixin =
  demoApp {
    values+: {
      replicaCount: 2,
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

local policy = latencyGradientPolicy({
  policyName: 'demo1-demo-app',
  serviceSelector+: {
    service: 'demo1-demo-app.demoapp.svc.cluster.local',
  },
}).policy;

{
  configMap:
    k.core.v1.configMap.new('policies')
    + k.core.v1.configMap.metadata.withLabels({ 'fluxninja.com/validate': 'true' })
    + k.core.v1.configMap.metadata.withNamespace('aperture-system')
    + k.core.v1.configMap.withData({
      'demo1-demo-app.yaml': std.manifestYamlDoc(policy, quote_keys=false),
    }),
  demoapp: demoappMixin,
}
