local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/postgresql',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'postgresql',
    applyStrategy: 'server',
  },
  data: mixins,
}
