{
  accept_percentage: import './panels/accept_percentage.libsonnet',
  accepted_concurrency: import './panels/accepted_concurrency.libsonnet',
  average_load_multiplier: import './panels/average_load_multiplier.libsonnet',
  incoming_concurrency: import './panels/incoming_concurrency.libsonnet',
  rate_limiter: import './panels/rate_limiter.libsonnet',
  signal_average: import './panels/signal_average.libsonnet',
  signal_frequency_invalid: import './panels/signal_frequency_invalid.libsonnet',
  signal_frequency_valid: import './panels/signal_frequency_valid.libsonnet',
  throughput: import './panels/throughput.libsonnet',
  token_bucket_available_tokens: import './panels/token_bucket_available_tokens.libsonnet',
  token_bucket_capacity: import './panels/token_bucket_capacity.libsonnet',
  token_bucket_fillrate: import './panels/token_bucket_fillrate.libsonnet',
  wfq_scheduler_flows: import './panels/wfq_scheduler_flows.libsonnet',
  wfq_scheduler_heap_requests: import './panels/wfq_scheduler_heap_requests.libsonnet',
  workload_decisions_accepted: import './panels/workload_decisions_accepted.libsonnet',
  workload_decision_rejected: import './panels/workload_decision_rejected.libsonnet',
  workload_latency: import './panels/workload_latency.libsonnet',
  query: import './panels/query.libsonnet',

  // Grouped panels
  auto_scale: import './panels/auto_scale/auto_scale.libsonnet',
  signals: import './panels/signals/signals.libsonnet',
}
