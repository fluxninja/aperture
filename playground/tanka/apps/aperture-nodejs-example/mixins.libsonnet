local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-js-example',
    },
  },
  environment+:: {
    namespace: 'aperture-js-example',
    name: 'aperture-js-example',
  },
}
