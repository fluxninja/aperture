{
  accept_percentage: import './panels/accept_percentage.libsonnet',
  accepted_token_rate: import './panels/accepted_token_rate.libsonnet',
  average_load_multiplier: import './panels/average_load_multiplier.libsonnet',
  incoming_token_rate: import './panels/incoming_token_rate.libsonnet',
  rate_limiter: import './panels/rate_limiter.libsonnet',
  signal_average: import './panels/signal_average.libsonnet',
  signal_frequency: import './panels/signal_frequency.libsonnet',
  throughput: import './panels/throughput.libsonnet',
  token_bucket_available_tokens: import './panels/token_bucket_available_tokens.libsonnet',
  token_bucket_capacity: import './panels/token_bucket_capacity.libsonnet',
  token_bucket_fillrate: import './panels/token_bucket_fillrate.libsonnet',
  wfq_scheduler_flows: import './panels/wfq_scheduler_flows.libsonnet',
  wfq_scheduler_heap_requests: import './panels/wfq_scheduler_heap_requests.libsonnet',
  workload_decisions_accepted: import './panels/workload_decisions_accepted.libsonnet',
  workload_decision_rejected: import './panels/workload_decisions_rejected.libsonnet',
  workload_latency: import './panels/workload_latency.libsonnet',
  query: import './panels/query.libsonnet',
  quota_checks: import './panels/quota_checks.libsonnet',

  // Grouped panels
  adaptive_load_scheduler: import './panels/grouped/adaptive_load_scheduler.libsonnet',
  auto_scale: import './panels/grouped/auto_scale.libsonnet',
  signals: import './panels/grouped/signals.libsonnet',
  load_ramp: import './panels/grouped/load_ramp.libsonnet',
  quota_scheduler: import './panels/grouped/quota_scheduler.libsonnet',
  pgsql: import './panels/grouped/pgsql.libsonnet',

  // PGSQL panels
  pgsql_db_count: import './panels/pgsql/db_count.libsonnet',
  pgsql_max_connections: import './panels/pgsql/max_connections.libsonnet',
  pgsql_total_commits: import './panels/pgsql/total_commits.libsonnet',
  pgsql_total_operations: import './panels/pgsql/total_operations.libsonnet',
  pgsql_total_rollbacks: import './panels/pgsql/total_rollbacks.libsonnet',
  pgsql_backends: import './panels/pgsql/backends.libsonnet',

}
