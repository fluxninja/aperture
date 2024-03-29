local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';
local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local demoappValues = if std.objectHas(values, 'demoapp') then values.demoapp else {};

local istioInjectLabels = {
  extraLabels+: {
    'sidecar.istio.io/inject': 'true',
  },
  replicaCount: 2,
};

local commonValues = if std.objectHas(demoappValues, 'common') then demoappValues.common else {};

local application = {
  environment:: {
    namespace: 'demoapp',
  },
  values:: commonValues + istioInjectLabels,
  service1:
    helm.template('service1', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: std.mergePatch($.values, if std.objectHas(demoappValues, 'service1') then demoappValues.service1 else {}),
    }),
  service2:
    helm.template('service2', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: std.mergePatch($.values, if std.objectHas(demoappValues, 'service2') then demoappValues.service2 else {}),
    }),
  service3:
    helm.template('service3', 'charts/demo-app', {
      namespace: $.environment.namespace,
      values: std.mergePatch($.values, if std.objectHas(demoappValues, 'service3') then demoappValues.service3 else {}),
    }),
};

application
