local workload_accepted = import '../workload_decisions_accepted.libsonnet';
local workload_rejected = import '../workload_decisions_rejected.libsonnet';
local workload_latency = import '../workload_latency.libsonnet';
local incoming_concurrency = import '../incoming_concurrency.libsonnet';
local accepted_concurrency = import '../accepted_concurrency.libsonnet';
local wfq_scheduler_flows = import '../wfq_scheduler_flows.libsonnet';
local wfq_scheduler_heap_requests = import '../wfq_scheduler_heap_requests.libsonnet';

function(cfg) {
  panels: [
    workload_accepted(cfg).panel,
    workload_rejected(cfg).panel,
    workload_latency(cfg).panel,
    incoming_concurrency(cfg).panel,
    accepted_concurrency(cfg).panel,
    wfq_scheduler_flows(cfg).panel,
    wfq_scheduler_heap_requests(cfg).panel,
  ],
}
