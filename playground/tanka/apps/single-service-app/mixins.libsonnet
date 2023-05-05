local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';

local helpers = import 'ninja/helpers.libsonnet';
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local demoappValues = if std.objectHas(values, 'demoapp') then values.demoapp else {};

local helm = tanka.helm.new(helpers.helmChartsRoot);

local application = {
  environment:: {
    namespace: 'demoapp',
  },
  values:: if std.objectHas(demoappValues, 'common') then demoappValues.common else {},
  service1:
    helm.template('service1', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: std.mergePatch($.values, if std.objectHas(demoappValues, 'service1') then demoappValues.service1 else {}),
    }),
};

application
