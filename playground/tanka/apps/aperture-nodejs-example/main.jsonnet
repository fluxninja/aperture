local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-nodejs-example',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-nodejs-example',
    applyStrategy: 'server',
  },
  data: mixins,
}
