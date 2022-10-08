local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-go-example',
    },
  },
  environment+:: {
    namespace: 'aperture-go-example',
    name: 'aperture-go-example',
  },
}
