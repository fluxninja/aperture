local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/grafonnet-v9.4.0/main.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  // Checkpoint Comparison (Requested vs Scheduled)
  local checkpointComparisonTargets = [
    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,infra_meter_name="%(infra_meter)s",type="requested"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Requested'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,infra_meter_name="%(infra_meter)s",type="scheduled"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Scheduled'),
  ],

  local checkpointComparison = timeSeriesPanel('Checkpoint Comparison', datasource, 'checkpoints/sec', stringFilters, targets=checkpointComparisonTargets),

  // Commits vs Rollbacks
  local commitVsRollbackTargets = [
    g.query.prometheus.new(datasource, 'rate(postgresql_commits_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval])' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Commits'),

    g.query.prometheus.new(datasource, 'rate(postgresql_rollbacks_total{%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval])' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rollbacks'),
  ],

  local commitVsRollback = timeSeriesPanel('Commits vs Rollbacks', datasource, 'transactions/sec', stringFilters, targets=commitVsRollbackTargets),

  // Block Reads (Heap vs Index)
  local blockReadsTargets = [
    g.query.prometheus.new(datasource, 'sum(rate(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_read"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Heap Read'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="idx_read"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Index Read'),
  ],

  local blockReads = timeSeriesPanel('Block Reads', datasource, 'blocks/sec', stringFilters, targets=blockReadsTargets),

  // Operations (Insert, Delete, Update, Hot Update)
  local operationsTargets = [
    g.query.prometheus.new(datasource, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="ins"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Insert'),

    g.query.prometheus.new(datasource, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="del"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Delete'),

    g.query.prometheus.new(datasource, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="upd"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Update'),

    g.query.prometheus.new(datasource, 'sum by (operation) (rate(postgresql_operations_total{%(filters)s,infra_meter_name="%(infra_meter)s",operation="hot_upd"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Hot Update'),
  ],

  local operations = timeSeriesPanel('Database Operations', datasource, 'operations/sec', stringFilters, targets=operationsTargets),

  local bufferWritesTargets = [
    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_buffers_writes_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="backend"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Backend'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_buffers_writes_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="backend_fsync"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Backend Fsync'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_buffers_writes_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="bgwriter"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Bgwriter'),

    g.query.prometheus.new(datasource, 'sum(rate(postgresql_bgwriter_buffers_writes_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="checkpoints"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Checkpoints'),
  ],


  local bufferWrites = timeSeriesPanel('Buffer Writes', datasource, 'writes/sec', stringFilters, targets=bufferWritesTargets),

  checkpointComparison: checkpointComparison.panel,
  commitVsRollback: commitVsRollback.panel,
  blockReads: blockReads.panel,
  operations: operations.panel,
  bufferWrites: bufferWrites.panel,
}
