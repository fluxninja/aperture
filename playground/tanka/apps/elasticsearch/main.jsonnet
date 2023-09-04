local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/elasticsearch',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'elasticsearch',
    applyStrategy: 'server',
  },
  data: mixins,
}
