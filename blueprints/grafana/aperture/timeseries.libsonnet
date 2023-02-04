local lib = import '../grafana.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local prometheus = grafana.prometheus;
local timeSeriesPanel = lib.TimeSeriesPanel;


{
  new(title, datasource, axisLabel=null, unit=null)::
    timeSeriesPanel.new(
      title=title,
      datasource=datasource,
      span=24,
      min_span=24,
      axis_label=axisLabel,
    ) +
    timeSeriesPanel.withOptions(
      timeSeriesPanel.options.withLegend() +
      timeSeriesPanel.options.withTooltip()
    ) +
    timeSeriesPanel.withDefaults(
      timeSeriesPanel.defaults.withColorMode('palette-classic') +
      timeSeriesPanel.defaults.withCustom({
        drawStyle: 'line',
        lineInterpolation: 'linear',
        lineWidth: 1,
        pointSize: 5,
        scaleDistribution: {
          type: 'linear',
        },
        showPoints: 'auto',
        spanNulls: false,
      })
    ) +
    timeSeriesPanel.withFieldConfig(
      timeSeriesPanel.fieldConfig.withDefaults(
        timeSeriesPanel.fieldConfig.defaults.withColor({
          mode: 'palette-classic',
        }) +
        timeSeriesPanel.fieldConfig.defaults.withUnit(unit) +
        timeSeriesPanel.fieldConfig.defaults.withCustom({
          [if axisLabel != null then 'axisLabel']: axisLabel,
          axisPlacement: 'auto',
          barAlignment: 0,
          fillOpacity: 0,
          gradientMode: 'none',
          drawStyle: 'line',
          hideFrom: {
            legend: false,
            tooltip: false,
            viz: false,
          },
          lineInterpolation: 'linear',
          lineWidth: 1,
          pointSize: 5,
          scaleDistribution: {
            type: 'linear',
          },
          showPoints: 'auto',
          spanNulls: false,
          stacking: {
            group: 'A',
            mode: 'none',
          },
          thresholdsStyle: {
            mode: 'off',
          },
        })
      )
    ),
}
