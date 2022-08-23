local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'environment/tilt/apps/bookinfo',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'bookinfo',
    applyStrategy: 'server',
  },
  data: mixins,
}
