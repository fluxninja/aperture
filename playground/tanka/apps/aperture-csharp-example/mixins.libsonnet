local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

sdk
{
  values+:: {
    image+: {
      repository: 'docker.io/fluxninja/aperture-csharp-example',
    },
  },
  environment+:: {
    namespace: 'aperture-csharp-example',
    name: 'aperture-csharp-example',
  },
}
