local utils = import '../../policy-utils.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

local newGraphPanel(title, datasource, query, axisLabel='', unit='') =
  graphPanel.new(
    title=title,
    datasource=datasource,
    labelY1=axisLabel,
    formatY1=unit,
  )
  .addTarget(
    prometheus.target(
      expr=query,
      intervalFactor=1,
    )
  );


local dashboardWithPanels(dashboardParams, filters) =
  local datasource = dashboardParams.datasource;
  local dsName = datasource.name;

  local FluxMeterPanel = newGraphPanel('FluxMeter latency Query', dsName, 'sum(increase(flux_meter_sum{%(filters)s}[$__rate_interval])) / sum(increase(flux_meter_count{%(filters)s}[$__rate_interval]))' % { filters: filters }, 'Latency (ms)', 'ms');

  local GcDurationPanel = newGraphPanel('Garbage Collector duration', dsName, 'avg(java_lang_Copy_LastGcInfo_duration{k8s_pod_name=~"service3-demo-app-.*"})');

  local CpuUsagePanel = newGraphPanel('CPU usage', dsName, 'avg(java_lang_OperatingSystem_CpuLoad{k8s_pod_name=~"service3-demo-app-.*"})');

  dashboard.new(
    title='Aperture Service Protection',
    editable=true,
    schemaVersion=18,
    refresh=dashboardParams.refresh_interval,
    time_from=dashboardParams.time_from,
    time_to=dashboardParams.time_to
  )
  .addTemplate(
    {
      current: {
        text: 'default',
        value: dsName,
      },
      hide: 0,
      label: 'Data Source',
      name: 'datasource',
      options: [],
      query: 'prometheus',
      refresh: 1,
      regex: datasource.filter_regex,
      type: 'datasource',
    }
  )
  .addPanel(FluxMeterPanel, gridPos={ h: 10, w: 24, x: 0, y: 10 })
  .addPanel(GcDurationPanel, gridPos={ h: 10, w: 24, x: 0, y: 20 })
  .addPanel(CpuUsagePanel, gridPos={ h: 10, w: 24, x: 0, y: 30 });

function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,
  local filters = utils.dictToPrometheusFilter(params.dashboard.extra_filters { flux_meter_name: policyName, policy_name: policyName, flow_status: 'OK' }),

  local dashboardDef = dashboardWithPanels(params.dashboard, filters),

  dashboard: dashboardDef,
}
