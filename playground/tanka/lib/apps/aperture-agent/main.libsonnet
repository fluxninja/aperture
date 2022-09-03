local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'aperture-system',
    includeCrds: true,
    name: 'agent',
  },
  values:: {
  },
  'aperture-agent':
    helm.template($.environment.name, 'charts/aperture-agent', {
      namespace: $.environment.namespace,
      includeCrds: $.environment.includeCrds,
      values: $.values,
    }),
}
