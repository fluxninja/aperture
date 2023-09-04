local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'elasticsearch',
  },
  values:: {
    global+: {
      kibanaEnabled: true,
    },
  },
  elasticsearch:
    helm.template('elasticsearch', 'charts/elasticsearch', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values,
    }),
}
