local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'environment/tilt/apps/demoapp',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'demoapp',
    applyStrategy: 'server',
  },
  data: mixins,
}
