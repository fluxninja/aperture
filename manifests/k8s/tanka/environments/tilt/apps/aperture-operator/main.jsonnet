local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'environment/tilt/apps/aperture-operator',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-system',
    applyStrategy: 'server',
  },
  data: mixins,
}
