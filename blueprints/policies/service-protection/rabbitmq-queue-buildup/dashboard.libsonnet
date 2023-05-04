local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local graphPanel = grafana.graphPanel;

function(cfg) {
  local params = config.common + config.dashboard + cfg,
  local policyName = params.policy_name,
  local queueName = params.queue_name,
  local ds = params.datasource,
  local dsName = ds.name,
  local refresh = params.refresh_interval,
  local time_from = params.time_from,
  local time_to = params.time_to,

  local queueBuildupPanel =
    graphPanel.new(
      title='Queue Buildup',
      datasource=dsName,
    )
    .addTarget(
      prometheus.target(
        expr=(
          'sum(rabbitmq_message_current{rabbitmq_queue_name="%(queue_name)s",state="ready"})' % { queue_name: queueName }
        ),
        intervalFactor=1,
      ),
    ),

  local workloadDecisionsAccepted =
    graphPanel.new(
      title='Workload Decisions (accepted)',
      datasource=dsName,
      format='reqps',
    )
    .addTarget(
      prometheus.target(
        expr=(
          'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policy_name)s",decision_type="DECISION_TYPE_ACCEPTED"}[$__rate_interval]))' % { policy_name: policyName }
        ),
        intervalFactor=1,
      ),
    ),

  local workloadDecisionsRejected =
    graphPanel.new(
      title='Workload Decisions (accepted)',
      datasource=dsName,
      format='reqps',
    )
    .addTarget(
      prometheus.target(
        expr=(
          'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policy_name)s",decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))' % { policy_name: policyName }
        ),
        intervalFactor=1,
      ),
    ),

  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      schemaVersion=36,
      editable=true,
      refresh=refresh,
      time_from=time_from,
      time_to=time_to,
    )
    .addTemplate(
      grafana.template.datasource(
        name='datasource',
        query='prometheus',
        label='Data Source',
        current=dsName,
        hide=0,
        regex=ds.filter_regex,
      )
    )
    .addPanel(
      panel=queueBuildupPanel,
      gridPos={ h: 10, w: 24, x: 0, y: 0 }
    )
    .addPanel(
      panel=workloadDecisionsAccepted,
      gridPos={ h: 10, w: 24, x: 0, y: 10 }
    )
    .addPanel(
      panel=workloadDecisionsRejected,
      gridPos={ h: 10, w: 24, x: 0, y: 20 }
    ),

  dashboard: dashboardDef,
}
