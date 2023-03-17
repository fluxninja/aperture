local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);

local enableNginx = std.extVar('ENABLE_DEMO_APP_NGINX') == 'true';

local application = {
  environment:: {
    namespace: 'demoapp',
  },
  values:: {
  },
  service1:
    helm.template('service1', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: $.values {
        nginx: {
          enabled: enableNginx,
          resources: {
            requests: {
              cpu: '1',
              memory: '1024Mi',
            },
            limits: {
              cpu: '4',
              memory: '4096Mi',
            },
          },
        },
      },
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
};

application
