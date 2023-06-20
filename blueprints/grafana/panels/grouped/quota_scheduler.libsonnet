local accepted_concurrency = import '../accepted_concurrency.libsonnet';
local incoming_concurrency = import '../incoming_concurrency.libsonnet';
local quota_checks = import '../quota_checks.libsonnet';
local wfq_scheduler_flows = import '../wfq_scheduler_flows.libsonnet';
local wfq_scheduler_heap_requests = import '../wfq_scheduler_heap_requests.libsonnet';
local workload_accepted = import '../workload_decisions_accepted.libsonnet';
local workload_rejected = import '../workload_decisions_rejected.libsonnet';
local workload_latency = import '../workload_latency.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    quota_checks(cfg).panel,
    workload_accepted(cfg).panel,
    workload_rejected(cfg).panel,
    workload_latency(cfg).panel,
    incoming_concurrency(cfg).panel,
    accepted_concurrency(cfg).panel
    + g.panel.timeSeries.gridPos.withX(12),
    wfq_scheduler_flows(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12),
    wfq_scheduler_heap_requests(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.gridPos.withX(12),
  ],
}
