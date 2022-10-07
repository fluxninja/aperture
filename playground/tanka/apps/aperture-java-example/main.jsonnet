local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-java-example',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-java-example',
    applyStrategy: 'server',
  },
  data: mixins,
}
