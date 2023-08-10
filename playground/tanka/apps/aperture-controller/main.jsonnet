local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-controller',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'default',
    applyStrategy: 'server',
  },
  data: mixins,
}
