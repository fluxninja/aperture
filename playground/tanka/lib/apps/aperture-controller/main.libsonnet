local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'aperture-controller',
    includeCrds: true,
    name: 'controller',
  },
  values:: {
  },
  'aperture-controller':
    helm.template($.environment.name, 'charts/aperture-controller', {
      namespace: $.environment.namespace,
      includeCrds: $.environment.includeCrds,
      values: $.values,
    }),
}
