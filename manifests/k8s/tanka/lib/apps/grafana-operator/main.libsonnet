local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local certmanager = import 'ninja/certmanager.libsonnet';
local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: error 'environment.namespace must be set',
    includeCrds: true,
    name: 'grafana-operator',
  },
  values:: {},
  operator:
    helm.template($.environment.name, 'charts/grafana-operator', {
      namespace: $.environment.namespace,
      includeCrds: $.environment.includeCrds,
      values: $.values,
    }),
}
