local application = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/demoui',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'demoui',
    applyStrategy: 'server',
  },
  data: application,
}
