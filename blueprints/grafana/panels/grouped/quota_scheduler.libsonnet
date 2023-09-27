local accepted_token_rate = import '../accepted_token_rate.libsonnet';
local incoming_token_rate = import '../incoming_token_rate.libsonnet';
local request_in_queue_duration = import '../request_in_queue_duration.libsonnet';
local wfq_scheduler_flows = import '../wfq_scheduler_flows.libsonnet';
local wfq_scheduler_heap_requests = import '../wfq_scheduler_heap_requests.libsonnet';
local workload_decisions = import '../workload_decisions.libsonnet';
local workload_accepted = import '../workload_decisions_accepted.libsonnet';
local workload_rejected = import '../workload_decisions_rejected.libsonnet';
local workload_latency = import '../workload_latency.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    workload_decisions(cfg).panel
    + g.panel.timeSeries.gridPos.withY(10),
    workload_accepted(cfg).panel
    + g.panel.timeSeries.gridPos.withY(20),
    workload_rejected(cfg).panel
    + g.panel.timeSeries.gridPos.withY(30),
    workload_latency(cfg).panel
    + g.panel.timeSeries.gridPos.withY(40),
    request_in_queue_duration(cfg).panel
    + g.panel.timeSeries.gridPos.withY(50),
    incoming_token_rate(cfg).panel
    + g.panel.timeSeries.gridPos.withY(60),
    accepted_token_rate(cfg).panel
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(60),
    wfq_scheduler_flows(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.timeSeries.gridPos.withY(70),
    wfq_scheduler_heap_requests(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(70),
  ],
}
