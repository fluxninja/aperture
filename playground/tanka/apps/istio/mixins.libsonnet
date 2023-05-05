local istioApp = import 'apps/istio/main.libsonnet';
local valuesStr = std.extVar('VALUES');
local values = if valuesStr != '' then std.parseYaml(valuesStr) else {};
local istioValues = if std.objectHas(values, 'istio') then values.istio else {};

local istioAppMixin =
  istioApp {
    values+: std.mergePatch({
      base+: {},
      istiod+: {},
      gateway+: {},
      envoyfilter+: {},
    }, istioValues),
  };

istioAppMixin
