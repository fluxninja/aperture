---
title: Agent
sidebar_position: 2
sidebar_label: Agent
---

## FluxMeter

### Metrics

| Name       | Type      | Labels                                                                                             | Unit | Description              |
| ---------- | --------- | -------------------------------------------------------------------------------------------------- | ---- | ------------------------ |
| flux_meter | Histogram | flux_meter_name, decision_type, response_status, http_status_code, feature_status, attribute_found | ms   | Flow's workload duration |

### Labels

| Name             | Example | Description |
| ---------------- | ------- | ----------- |
| flux_meter_name  |         |             |
| decision_type    |         |             |
| response_status  |         |             |
| http_status_code |         |             |
| feature_status   |         |             |
| attribute_found  |         |             |

## ConcurrencyLimiter

### Metrics

| Name                    | Type    | Labels                                                    | Unit            | Description                                               |
| ----------------------- | ------- | --------------------------------------------------------- | --------------- | --------------------------------------------------------- |
| workload_latency_ms     | Summary | policy_name, policy_hash, component_index, workload_index | ms              | Latency summary of workload                               |
| workload_requests_total | Counter | policy_name, policy_hash, component_index, workload_index | count (no unit) | A counter of workload requests                            |
| incoming_concurrency_ms | Counter | policy_name, policy_hash, component_index                 | ms              | A counter measuring incoming concurrency into Scheduler   |
| accepted_concurrency_ms | Counter | policy_name, policy_hash, component_index                 | ms              | A counter measuring the concurrency admitted by Scheduler |

### Labels

| Name            | Example | Description |
| --------------- | ------- | ----------- |
| policy_name     |         |             |
| policy_hash     |         |             |
| workload_index  |         |             |
| component_index |         |             |

## RateLimiter

### Metrics

| Name                 | Type    | Labels                                    | Unit            | Description                                                        |
| -------------------- | ------- | ----------------------------------------- | --------------- | ------------------------------------------------------------------ |
| rate_limiter_counter | Counter | policy_name, policy_hash, component_index | count (no unit) | A counter measuring the number of times Rate Limiter was triggered |

### Labels

| Name            | Example | Description |
| --------------- | ------- | ----------- |
| policy_name     |         |             |
| policy_hash     |         |             |
| component_index |         |             |

## Classifier

### Metrics

| Name               | Type    | Labels                                     | Unit            | Description                                                      |
| ------------------ | ------- | ------------------------------------------ | --------------- | ---------------------------------------------------------------- |
| classifier_counter | Counter | policy_name, policy_hash, classifier_index | count (no unit) | A counter measuring the number of times classifier was triggered |

### Labels

| Name             | Example | Description |
| ---------------- | ------- | ----------- |
| policy_name      |         |             |
| policy_hash      |         |             |
| classifier_index |         |             |

## Dataplane Summary

### FlowControl Metrics

| Name                             | Type    | Labels        | Unit            | Description                                     |
| -------------------------------- | ------- | ------------- | --------------- | ----------------------------------------------- |
| flowcontrol_requests_total       | Counter |               | count (no unit) | Total number of aperture check requests handled |
| flowcontrol_decisions_total      | Counter | decision_type | count (no unit) | Number of aperture check decisions              |
| flowcontrol_error_reasons_total  | Counter | error_reason  | count(no unit)  | Number of error reasons other than unspecified  |
| flowcontrol_reject_reasons_total | Counter | reject_reason | count (no unit) | Number of reject reasons other than unspecified |

### FlowControl Labels

| Name          | Example | Description |
| ------------- | ------- | ----------- |
| decision_type |         |             |
| error_reason  |         |             |
| reject_reason |         |             |

### Distributed Cache Metrics

| Name                    | Type  | Labels                                     | Unit            | Description                                                                           |
| ----------------------- | ----- | ------------------------------------------ | --------------- | ------------------------------------------------------------------------------------- |
| distcache_entries_total | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Total number of entries in the DMap                                                   |
| distcache_delete_hits   | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Number of deletion requests resulting in an item being removed in the DMap            |
| distcache_delete_misses | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Number of deletion requests for missing keys in the DMap                              |
| distcache_get_misses    | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Number of entries that have been requested and not found in the DMap                  |
| distcache_get_hits      | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Number of entries that have been requested and found present in the DMap              |
| distcache_evicted_total | Gauge | distcache_member_id, distcache_member_name | count (no unit) | Total number of entries removed from cache to free memory for new entries in the DMap |

### Distributed Cache Labels

| Name                  | Example | Description |
| --------------------- | ------- | ----------- |
| distcache_member_id   |         |             |
| distcache_member_name |         |             |

### Scheduler Metrics

| Name                                | Type  | Labels                                    | Unit            | Description                                                           |
| ----------------------------------- | ----- | ----------------------------------------- | --------------- | --------------------------------------------------------------------- |
| wfq_flows_total                     | Gauge | policy_name, policy_hash, component_index | count (no unit) | A gauge that tracks the number of flows in the WFQScheduler           |
| wfq_requests_total                  | Gauge | policy_name, policy_hash, component_index | count (no unit) | A gauge that tracks the number of queued requests in the WFQScheduler |
| token_bucket_lm_ratio               | Gauge | policy_name, policy_hash, component_index | percentage      | A gauge that tracks the load multiplier                               |
| token_bucket_fill_rate              | Gauge | policy_name, policy_hash, component_index | tokens/s        | A gauge that tracks the fill rate of token bucket                     |
| token_bucket_capacity_total         | Gauge | policy_name, policy_hash, component_index | count (no unit) | A gauge that tracks the capacity of token bucket                      |
| token_bucket_available_tokens_total | Gauge | policy_name, policy_hash, component_index | count (no unit) | A gauge that tracks the number of tokens available in token bucket    |

### Scheduler Labels

| Name            | Example | Description |
| --------------- | ------- | ----------- |
| policy_name     |         |             |
| policy_hash     |         |             |
| component_index |         |             |
