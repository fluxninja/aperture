{
  accept_percentage: import 'accept_percentage.libsonnet',
  accepted_concurrency: import 'accepted_concurrency.libsonnet',
  average_load_multiplier: import 'average_load_multiplier.libsonnet',
  incoming_concurrency: import 'incoming_concurrency.libsonnet',
  rate_limiter: import 'rate_limiter.libsonnet',
  signal_average: import 'signal_average.libsonnet',
  signal_frequency_invalid: import 'signal_frequency_invalid.libsonnet',
  signal_frequency_valid: import 'signal_frequency_valid.libsonnet',
  throughput: import 'throughput.libsonnet',
  token_bucket_available_tokens: import 'token_bucket_available_tokens.libsonnet',
  token_bucket_capacity: import 'token_bucket_capacity.libsonnet',
  token_bucket_fillrate: import 'token_bucket_fillrate.libsonnet',
  wfq_scheduler_flows: import 'wfq_scheduler_flows.libsonnet',
  wfq_scheduler_heap_requests: import 'wfq_scheduler_heap_requests.libsonnet',
  workload_decisions_accepted: import 'workload_decisions_accepted.libsonnet',
  workload_decision_rejected: import 'workload_decision_rejected.libsonnet',
  workload_latency: import 'workload_latency.libsonnet',
  query: import 'query.libsonnet',

  // Grouped panels
  auto_scale: import './auto_scale/auto_scale.libsonnet',
}
