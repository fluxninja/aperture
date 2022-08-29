local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
  },
  demo1:
    helm.template('demo1', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
}
