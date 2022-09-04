local mixins = import 'mixins.libsonnet';

function(apiServer='API SERVER MISSING') {
  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'apps/aperture-agent',
  },
  spec: {
    apiServer: apiServer,
    namespace: 'aperture-agent',
    applyStrategy: 'server',
  },
  data: mixins,
}
