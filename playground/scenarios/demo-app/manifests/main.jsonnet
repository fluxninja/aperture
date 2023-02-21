local application = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/demoapp',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'demoapp',
    applyStrategy: 'server',
  },
  data: application,
}
