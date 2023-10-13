local accepted_token_rate = import '../accepted_token_rate.libsonnet';
local avg_preemption_chart = import '../avg_preemption_chart.libsonnet';
local avg_preemption_time_series = import '../avg_preemption_time_series.libsonnet';
local incoming_token_rate = import '../incoming_token_rate.libsonnet';
local request_in_queue_duration = import '../request_in_queue_duration.libsonnet';
local request_queue_duration_bar = import '../request_queue_duration_bar.libsonnet';
local total_accepted_requests = import '../total_accepted_requests.libsonnet';
local total_accepted_tokens = import '../total_accepted_tokens.libsonnet';
local total_incoming_tokens = import '../total_incoming_tokens.libsonnet';
local total_rejected_requests = import '../total_rejected_requests.libsonnet';
local total_rejected_tokens = import '../total_rejected_tokens.libsonnet';
local total_requests = import '../total_requests.libsonnet';
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
    total_requests(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(40),
    total_accepted_requests(cfg).panel
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(40),
    total_rejected_requests(cfg).panel
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(40),
    workload_latency(cfg).panel
    + g.panel.timeSeries.gridPos.withY(50),
    request_in_queue_duration(cfg).panel
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(60)
    + g.panel.timeSeries.gridPos.withW(12),
    request_queue_duration_bar(cfg).panel
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.barGauge.gridPos.withY(60)
    + g.panel.barGauge.gridPos.withW(12),
    avg_preemption_time_series(cfg).panel
    + g.panel.timeSeries.gridPos.withY(70)
    + g.panel.timeSeries.gridPos.withW(12),
    avg_preemption_chart(cfg).panel
    + g.panel.barChart.gridPos.withX(12)
    + g.panel.barChart.gridPos.withY(70)
    + g.panel.barChart.gridPos.withW(12),
    incoming_token_rate(cfg).panel
    + g.panel.timeSeries.gridPos.withY(80),
    accepted_token_rate(cfg).panel
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(80),
    total_incoming_tokens(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(90),
    total_accepted_tokens(cfg).panel
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(90),
    total_rejected_tokens(cfg).panel
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(90),
    wfq_scheduler_flows(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.timeSeries.gridPos.withY(100),
    wfq_scheduler_heap_requests(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(100),
  ],
}
