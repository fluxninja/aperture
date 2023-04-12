local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

{
  environment:: {
    namespace: 'istio-system',
  },
  values:: {
    base: {},
    istiod: {},
    envoyfilter: {
      enableAuthzRequestBodyBuffering: true,
    },
  },
  base:
    helm.template('base', 'charts/base', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values.base,
    }),
  istiod:
    helm.template('istiod', 'charts/istiod', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values.istiod,
    }),
  gateway:
    helm.template('gateway', 'charts/gateway', {
      namespace: $.environment.namespace,
      includeCrds: true,
      values: $.values.gateway,
    }),
  envoyfilter:
    helm.template('envoyfilter', 'charts/istioconfig', {
      namespace: $.environment.namespace,
      values: $.values.envoyfilter,
    }),
}
