local graphqlDemoApp = import './main.libsonnet';

local k = import 'github.com/jsonnet-libs/k8s-libsonnet/1.22/main.libsonnet';

local wavepoolConfigName = 'wavepool-graphql-config';
local wavePoolGenerator = 'wavepool-graphql-generator';

local graphQLDemoappMixin =
  graphqlDemoApp {
    values+: {
      replicaCount: 2,
      simplesrv+: {
        image: {
          repository: 'docker.io/fluxninja/graphql-demo-app',
          tag: 'test',
        },
      },
      resources+: {
        limits+: {
          cpu: '100m',
          memory: '128Mi',
        },
        requests+: {
          cpu: '100m',
          memory: '128Mi',
        },
      },
    },
  };

local wavePoolConfigMixin =
  k.core.v1.configMap.new(wavepoolConfigName)
  + k.core.v1.configMap.withData({
    'graphql_wavepool_generator.js': importstr '../../../load_generator/scenarios/graphql_load_test.js',
  });

local wavePoolGeneratorMixin =
  k.apps.v1.deployment.new(wavePoolGenerator)
  + k.apps.v1.deployment.spec.selector.withMatchLabels({ 'app.kubernetes.io/component': wavePoolGenerator })
  + k.apps.v1.deployment.spec.template.metadata.withLabels({ 'app.kubernetes.io/component': wavePoolGenerator, 'sidecar.istio.io/inject': 'false' })
  + k.apps.v1.deployment.spec.template.spec.withContainers({
    name: wavePoolGenerator,
    image: 'loadimpact/k6:latest',
    imagePullPolicy: 'Always',
    command: ['/bin/sh', '-xc'],
    args: ['while true; do k6 run -v /tmp/graphql_wavepool_generator.js; done'],
    resources: {
      limits: {
        cpu: '1',
        memory: '2Gi',
      },
    },
    volumeMounts: [{
      mountPath: '/tmp',
      name: 'graphql-js-file',
    }],
  })
  + k.apps.v1.deployment.spec.template.spec.withVolumes({
    configMap: {
      name: wavepoolConfigName,
    },
    name: 'graphql-js-file',
  });

{
  demoapp: graphQLDemoappMixin,
  wavePoolconfigMap: wavePoolConfigMixin,
  wavepoolDeployment: wavePoolGeneratorMixin,
}
