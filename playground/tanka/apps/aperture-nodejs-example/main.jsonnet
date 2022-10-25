local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-js-example',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-js-example',
    applyStrategy: 'server',
  },
  data: mixins,
}
