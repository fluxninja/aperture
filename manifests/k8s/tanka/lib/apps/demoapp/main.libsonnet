local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
  },
  serviceA:
    helm.template('service-a', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  serviceB:
    helm.template('service-b', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  serviceC:
    helm.template('service-c', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
}
