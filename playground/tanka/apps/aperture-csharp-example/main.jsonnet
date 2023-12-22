local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-csharp-example',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-csharp-example',
    applyStrategy: 'server',
  },
  data: mixins,
}
