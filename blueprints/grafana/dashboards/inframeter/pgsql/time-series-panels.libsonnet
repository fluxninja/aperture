local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  // Checkpoint Comparison (Requested vs Scheduled)
  local checkpointComparisonTargets = [
    g.query.prometheus.new(datasource, 'sum by (type) (rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ type }}'),
  ],

  local checkpointComparison = timeSeriesPanel('Checkpoint Rate', datasource, axisLabel='checkpoints/sec', unit='cps', targets=checkpointComparisonTargets),

  // Commits vs Rollbacks
  local commitVsRollbackTargets = [
    g.query.prometheus.new(datasource, 'sum(rate(postgresql_commits_total{%(filters)s"}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Commits'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_rollbacks_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rollbacks'),
  ],

  local commitVsRollback = timeSeriesPanel('Commits vs Rollbacks', datasource, axisLabel='transactions/sec', unit='cps', targets=commitVsRollbackTargets),

  // Block Reads (Heap vs Index)
  local blockReadsTargets = [
    g.query.prometheus.new(datasource, 'sum by (source) (rate(postgresql_blocks_read_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ source }}'),
  ],

  local blockReads = timeSeriesPanel('Block Reads', datasource, axisLabel='blocks/sec', unit='cps', targets=blockReadsTargets),

  // Operations (Insert, Delete, Update, Hot Update)
  local operationsTargets = [
    g.query.prometheus.new(datasource, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ operation }}'),
  ],

  local operations = timeSeriesPanel('Database Operations', datasource, axisLabel='operations/sec', unit='ops', targets=operationsTargets),

  local bufferWritesTargets = [
    g.query.prometheus.new(datasource, 'sum by (source) (rate(postgresql_bgwriter_buffers_writes_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ source }}'),
  ],


  local bufferWrites = timeSeriesPanel('Buffer Writes', datasource, axisLabel='writes/sec', unit='wps', targets=bufferWritesTargets),

  checkpointComparison: checkpointComparison,
  commitVsRollback: commitVsRollback,
  blockReads: blockReads,
  operations: operations,
  bufferWrites: bufferWrites,
}
