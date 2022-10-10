local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-java-example',
    },
  },
  environment+:: {
    namespace: 'aperture-java-example',
    name: 'aperture-java-example',
  },
}
