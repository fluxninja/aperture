local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-python-example',
    },
  },
  environment+:: {
    namespace: 'aperture-python-example',
    name: 'aperture-python-example',
  },
}
