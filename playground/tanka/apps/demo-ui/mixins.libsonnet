local tanka = import 'github.com/grafana/jsonnet-libs/tanka-util/main.libsonnet';
local helpers = import 'ninja/helpers.libsonnet';

local helm = tanka.helm.new(helpers.helmChartsRoot);
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local demouiValues = if std.objectHas(values, 'demoui') then values.demoui else {};

local commonValues = if std.objectHas(demouiValues, 'common') then demouiValues.common else {};

local application = {
  environment:: {
    namespace: 'demoui',
  },
  values:: commonValues,
  service:
    helm.template('service', 'charts/demo-ui', {
      namespace: $.environment.namespace,
      values: std.mergePatch(
        $.values,
        if std.objectHas(demouiValues, 'service') then demouiValues.service else {}
      ),
    }),
};

application
