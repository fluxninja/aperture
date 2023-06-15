local utils = import '../../../policies/policy-utils.libsonnet';
local quotaSchedular = import '../quota-scheduler/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local statPanel = grafana.statPanel;


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

  local maxId = std.reverse(std.sort(basequotaSchedular.panels, keyF=function(panel) '%s' % panel.id))[0].id,
  local maxPanelYAxis = std.reverse(std.sort(basequotaSchedular.panels, keyF=function(panel) panel.gridPos.y))[0].gridPos.y,

  local finaldashboard =
    basequotaSchedular {
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
