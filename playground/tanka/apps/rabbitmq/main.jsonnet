local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/rabbitmq',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'rabbitmq',
    applyStrategy: 'server',
  },
  data: mixins,
}
