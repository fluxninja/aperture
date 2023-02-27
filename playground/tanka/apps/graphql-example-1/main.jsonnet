local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

local application = {
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
    replicaCount: 2,
    simplesrv+: {
      image: {
        repository: 'docker.io/fluxninja/graphql-demo-app',
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
  service:
    helm.template('service', 'charts/graphql-demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
};

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/graphql-demoapp',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'demoapp',
    applyStrategy: 'server',
  },
  data: application,
}
