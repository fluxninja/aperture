local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'tanka/lib/ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
  },
  service:
    helm.template('service', 'charts/graphql-demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
}
