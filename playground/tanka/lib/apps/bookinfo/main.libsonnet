local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'bookinfo',
    name: 'bookinfo',
  },
  values:: {
  },
  bookinfo:
    helm.template($.environment.name, 'charts/bookinfo', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
}
