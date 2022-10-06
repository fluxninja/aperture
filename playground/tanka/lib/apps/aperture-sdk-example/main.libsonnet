local sdk = import 'aperture-sdk-example.libsonnet';

{
  local tl = self,
  values:: {},
  environment:: {},

  sdk: sdk(values=tl.values, environment=tl.environment),
}
