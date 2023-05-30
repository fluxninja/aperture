local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

local application = {
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
    replicaCount: 2,
    podAnnotations: {
      'prometheus.io/scrape': 'true',
      'prometheus.io/port': '8087',
    },
  },
  service1:
    helm.template('service1', 'charts/java-demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  service2:
    helm.template('service2', 'charts/java-demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
  service3:
    helm.template('service3', 'charts/java-demo-app', {
      namespace: $.environment.namespace,
      values: $.values,
    }),
};

application
