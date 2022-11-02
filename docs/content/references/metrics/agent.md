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

| Name             | Example                                        | Description                                                                                                   |
| ---------------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| flux_meter_name  | service1-demo-app                              | Name of the FluxMeter                                                                                         |
| decision_type    | DECISION_TYPE_ACCEPTED, DECISION_TYPE_REJECTED | Whether the flow was accepted or not                                                                          |
| response_status  | Error, OK                                      | A common label to denote OK or Error across all protocols                                                     |
| http_status_code | 200, 503                                       | HTTP status code                                                                                              |
| feature_status   |                                                | Feature status                                                                                                |
| valid            | true, false                                    | Label for specifying if metric is valid. A metric may be invalid if attribute is not found in flow telemetry. |

## ConcurrencyLimiter

### Metrics

| Name                    | Type    | Labels                                                    | Unit            | Description                                               |
| ----------------------- | ------- | --------------------------------------------------------- | --------------- | --------------------------------------------------------- |
| workload_latency_ms     | Summary | policy_name, policy_hash, component_index, workload_index | ms              | Latency summary of workload                               |
| workload_requests_total | Counter | policy_name, policy_hash, component_index, workload_index | count (no unit) | A counter of workload requests                            |
| incoming_concurrency_ms | Counter | policy_name, policy_hash, component_index                 | ms              | A counter measuring incoming concurrency into Scheduler   |
| accepted_concurrency_ms | Counter | policy_name, policy_hash, component_index                 | ms              | A counter measuring the concurrency admitted by Scheduler |

### Labels

| Name            | Example                                      | Description                                                     |
| --------------- | -------------------------------------------- | --------------------------------------------------------------- |
| policy_name     | service1-demo-app                            | Name of the policy.                                             |
| policy_hash     | 5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo= | Hash of the policy used for checking integrity of the policy.   |
| workload_index  | 0, 1, 2, default                             | Index of the workload in order of specification in the policy.  |
| component_index | 13                                           | Index of the component in order of specification in the policy. |

## RateLimiter

### Metrics

| Name                 | Type    | Labels                                    | Unit            | Description                                                        |
| -------------------- | ------- | ----------------------------------------- | --------------- | ------------------------------------------------------------------ |
| rate_limiter_counter | Counter | policy_name, policy_hash, component_index | count (no unit) | A counter measuring the number of times Rate Limiter was triggered |

### Labels

| Name            | Example                                      | Description                                                     |
| --------------- | -------------------------------------------- | --------------------------------------------------------------- |
| policy_name     | service1-demo-app                            | Name of the policy.                                             |
| policy_hash     | 5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo= | Hash of the policy used for checking integrity of the policy.   |
| component_index | 13                                           | Index of the component in order of specification in the policy. |

## Classifier

### Metrics

| Name               | Type    | Labels                                     | Unit            | Description                                                      |
| ------------------ | ------- | ------------------------------------------ | --------------- | ---------------------------------------------------------------- |
| classifier_counter | Counter | policy_name, policy_hash, classifier_index | count (no unit) | A counter measuring the number of times classifier was triggered |

### Labels

| Name             | Example                                      | Description                                                      |
| ---------------- | -------------------------------------------- | ---------------------------------------------------------------- |
| policy_name      | service1-demo-app                            | Name of the policy.                                              |
| policy_hash      | 5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo= | Hash of the policy used for checking integrity of the policy.    |
| classifier_index | 0, 1                                         | Index of the classifier in order of specification in the policy. |

## Dataplane Summary

### FlowControl Metrics

| Name                             | Type    | Labels        | Unit            | Description                                     |
| -------------------------------- | ------- | ------------- | --------------- | ----------------------------------------------- |
| flowcontrol_requests_total       | Counter |               | count (no unit) | Total number of aperture check requests handled |
| flowcontrol_decisions_total      | Counter | decision_type | count (no unit) | Number of aperture check decisions              |
| flowcontrol_error_reasons_total  | Counter | error_reason  | count(no unit)  | Number of error reasons other than unspecified  |
| flowcontrol_reject_reasons_total | Counter | reject_reason | count (no unit) | Number of reject reasons other than unspecified |

### FlowControl Labels

| Name          | Example                                                                                                                                              | Description                                              |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| decision_type | DECISION_TYPE_ACCEPTED, DECISION_TYPE_REJECTED                                                                                                       | Whether the flow was accepted or not                     |
| error_reason  | ERROR_NONE, ERROR_MISSING_TRAFFIC_DIRECTION, ERROR_INVALID_TRAFFIC_DIRECTION, ERROR_CONVERT_TO_MAP_STRUCT, ERROR_CONVERT_TO_REGO_AST, ERROR_CLASSIFY | Error reason for FlowControl Check response.             |
| reject_reason | REJECT_REASON_NONE, REJECT_REASON_RATE_LIMITED, REJECT_REASON_CONCURRENCY_LIMITED                                                                    | Reason why FlowControl Check response rejected the flow. |

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

| Name                  | Example            | Description                                               |
| --------------------- | ------------------ | --------------------------------------------------------- |
| distcache_member_id   | 384313659919819706 | Internal ID of distributed cache cluster member.          |
| distcache_member_name | 10.244.1.20:3320   | Internal unique name of distributed cache cluster member. |

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

| Name            | Example                                      | Description                                                     |
| --------------- | -------------------------------------------- | --------------------------------------------------------------- |
| policy_name     | service1-demo-app                            | Name of the policy.                                             |
| policy_hash     | 5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo= | Hash of the policy used for checking integrity of the policy.   |
| component_index | 13                                           | Index of the component in order of specification in the policy. |
