local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'rabbitmq',
  },
  values:: {},
  rabbitmq:
    helm.template('rabbitmq', 'charts/rabbitmq', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values,
    }),
}
