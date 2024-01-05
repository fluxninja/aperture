local sdk = import 'apps/aperture-sdk-example/main.libsonnet';

local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local sdkValues = if std.objectHas(values, 'sdk') then values.sdk else {};

sdk
{
  values+:: sdkValues {
    image+: {
      repository: 'docker.io/fluxninja/aperture-go-example',
    },
  },
  environment+:: {
    namespace: 'aperture-go-example',
    name: 'aperture-go-example',
  },
}
