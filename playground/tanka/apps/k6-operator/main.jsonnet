local k6_operator = import 'apps/k6-operator/main.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/k6-operator',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'demoapp',
    applyStrategy: 'server',
  },
  data: k6_operator,
}
