local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local timeSeriesPanel = lib.TimeSeriesPanel;

local newTimeSeriesPanel(title, datasource, query, axisLabel='', unit='') =
  local thresholds =
    {
      mode: 'absolute',
      steps: [
        { color: 'green', value: null },
        { color: 'red', value: 80 },
      ],
    };
  local target =
    prometheus.target(query, intervalFactor=1)
    + { range: true, editorMode: 'code' };
  aperture.timeSeriesPanel.new(title, datasource, axisLabel, unit)
  + timeSeriesPanel.withTarget(target)
  + timeSeriesPanel.defaults.withThresholds(thresholds)
  + timeSeriesPanel.withFieldConfigMixin(
    timeSeriesPanel.fieldConfig.withDefaultsMixin(
      timeSeriesPanel.fieldConfig.defaults.withThresholds(thresholds)
    )
  ) + {
    interval: '1s',
  };

function(cfg) {
  local params = config.common + config.dashboard + cfg,
  local policyName = params.policy_name,
  local queueName = params.queue_name,
  local ds = params.datasource,
  local dsName = ds.name,

  local QueueBuildupPanel =
    newTimeSeriesPanel('Queue Buildup',
                       dsName,
                       |||
                         sum(rabbitmq_message_current{rabbitmq_queue_name="%(queue_name)s",state="ready"})
                       ||| % { queue_name: queueName },
                       'Ready Messages'),

  local WorkloadDecisionsAccepted =
    newTimeSeriesPanel('Workload Decisions (accepted)', dsName, 'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policy_name)s",decision_type="DECISION_TYPE_ACCEPTED"}[$__rate_interval]))' % { policy_name: policyName }, 'Decisions', 'reqps'),

  local WorkloadDecisionsRejected =
    newTimeSeriesPanel('Workload Decisions (rejected)', dsName, 'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policy_name)s",decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))' % { policy_name: policyName }, 'Decisions', 'reqps'),


  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      editable=true,
      schemaVersion=18,
      refresh=params.refresh_interval,
      time_from=params.time_from,
      time_to=params.time_to
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
        refres: 1,
        regex: ds.filter_regex,
        type: 'datasource',
      }
    )
    .addPanel(QueueBuildupPanel, gridPos={ h: 10, w: 24, x: 0, y: 0 })
    .addPanel(WorkloadDecisionsAccepted, gridPos={ h: 10, w: 24, x: 0, y: 10 })
    .addPanel(WorkloadDecisionsRejected, gridPos={ h: 10, w: 24, x: 0, y: 20 }),

  dashboard: dashboardDef,
}
