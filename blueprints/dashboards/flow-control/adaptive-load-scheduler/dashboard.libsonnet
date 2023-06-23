local utils = import '../../../policies/policy-utils.libsonnet';
local quotaSchedular = import '../quota-scheduler/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local statPanel = grafana.statPanel;
local graphPanel = grafana.graphPanel;


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

local newStatPanel(graphTitle, datasource, graphQuery) =
  local target =
    prometheus.target(graphQuery) +
    {
      legendFormat: '{{ instance }} - {{ policy_name }}',
      editorMode: 'code',
      range: true,
    };

  statPanel.new(
    title=graphTitle,
    datasource=datasource,
  ).addTarget(target) +
  {
    fieldConfig: {
      defaults: {
        color: {
          mode: 'thresholds',
        },
        mappings: [],
        thresholds: {
          mode: 'absolute',
          steps: [
            { color: 'green', value: null },
          ],
        },
      },
      overrides: [],
    },
    options: {
      colorMode: 'value',
      graphMode: 'none',
      justifyMode: 'center',
      orientation: 'horizontal',
      reduceOptions: {
        calcs: ['lastNotNull'],
        fields: '',
        values: false,
      },
      textMode: 'auto',
    },
  };

local dashboardWithPanels(dashboardParams, filters) =
  {
    local datasource = dashboardParams.datasource,
    local dsName = datasource.name,

    local TotalBucketLoadSchedFactor =
      newStatPanel('Average Load Multiplier', dsName, 'avg(token_bucket_lm_ratio{%(filters)s})' % { filters: filters }),

    local TokenBucketBucketCapacity =
      newStatPanel('Token Bucket Bucket Capacity', dsName, 'avg(token_bucket_capacity_total{%(filters)s})' % { filters: filters })
      + {
        options+: {
          orientation: 'auto',
        },
      },

    local TokenBucketBucketFillRate =
      newStatPanel('Token Bucket Bucket FillRate', dsName, 'avg(token_bucket_fill_rate{%(filters)s})' % { filters: filters }) +
      {
        options+: {
          orientation: 'auto',
        },
      },

    local TokenBucketAvailableTokens =
      newStatPanel('Token Bucket Available Tokens', dsName, 'avg(token_bucket_available_tokens_total{%(filters)s})' % { filters: filters }) +
      {
        options+: {
          orientation: 'auto',
        },
      },
    local dashboardDef =
      dashboard.new(
        title=dashboardParams.title,
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
          refres: 1,
          regex: datasource.filter_regex,
          type: 'datasource',
        }
      )
      .addPanel(TotalBucketLoadSchedFactor, gridPos={ h: 6, w: 6, x: 0, y: 50 })
      .addPanel(TokenBucketBucketCapacity, gridPos={ h: 6, w: 6, x: 6, y: 50 })
      .addPanel(TokenBucketBucketFillRate, gridPos={ h: 6, w: 6, x: 12, y: 50 })
      .addPanel(TokenBucketAvailableTokens, gridPos={ h: 6, w: 6, x: 18, y: 50 }),

    dashboard: dashboardDef,

  };

function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,
  local filters = utils.dictToPrometheusFilter(params.dashboard.extra_filters { policy_name: policyName }),

  local LoadSchedular = dashboardWithPanels(params.dashboard, filters).dashboard,

  local basequotaSchedular = quotaSchedular(params).dashboard,

  local overloadConfirmationPanels =
    if std.objectHas(params.policy, 'service_protection_core') && std.objectHas(params.policy.service_protection_core, 'overload_confirmations') then [
      local query = params.policy.service_protection_core.overload_confirmations[idx];
      newGraphPanel(
        'Overload Confirmation Query %s - %s - %0.3f' % [(idx + 1), query.operator, query.threshold],
        params.dashboard.datasource.name,
        params.policy.service_protection_core.overload_confirmations[idx].query_string,
      ) {
        id: idx + 1,
        gridPos: { x: 0, y: 10 + (idx * 10), w: 24, h: 10 },
      }
      for idx in std.range(0, std.length(params.policy.service_protection_core.overload_confirmations) - 1)
    ]
    else [],

  local overloadConfirmationPanelsLength = std.length(overloadConfirmationPanels),

  local quotaSchedularDashboard = basequotaSchedular {
    panels: overloadConfirmationPanels + [
      basequotaSchedular.panels[panel_idx] {
        id: overloadConfirmationPanelsLength + panel_idx + 1,
        gridPos+: { y: 10 + (10 * overloadConfirmationPanelsLength) + basequotaSchedular.panels[panel_idx].gridPos.y },
      }
      for panel_idx in std.range(0, std.length(basequotaSchedular.panels) - 1)
    ],
  },

  local maxId = std.reverse(std.sort(quotaSchedularDashboard.panels, keyF=function(panel) '%s' % panel.id))[0].id,
  local maxPanelYAxis = std.reverse(std.sort(quotaSchedularDashboard.panels, keyF=function(panel) panel.gridPos.y))[0].gridPos.y,

  local finaldashboard =
    quotaSchedularDashboard {
      panels+: [
        LoadSchedular.panels[panel_idx] {
          id: maxId + panel_idx + 1,
          gridPos: { y: maxPanelYAxis + 10, x: 6 * panel_idx, w: 6, h: 6 },
        }
        for panel_idx in std.range(0, std.length(LoadSchedular.panels) - 1)
      ],
    },

  dashboard: finaldashboard,
}
