local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-nodejs-example',
    },
  },
  environment+:: {
    namespace: 'aperture-nodejs-example',
    name: 'aperture-nodejs-example',
  },
}
