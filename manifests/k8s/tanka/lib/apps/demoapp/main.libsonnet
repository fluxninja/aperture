local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
  },
  service1:
    helm.template('service1', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  service2:
    helm.template('service2', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  service3:
    helm.template('service3', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
}
