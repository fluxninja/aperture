local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'istio-system',
  },
  values:: {
    istio: {},
    envoyfilter: {},
  },
  istio:
    helm.template('istio', 'charts/istio', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values.istio,
    }),
  envoyfilter:
    helm.template('envoyfilter', 'charts/istioconfig', {
      namespace: $.environment.namespace,
      values: $.values.envoyfilter,
    }),
}
