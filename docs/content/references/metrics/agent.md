---
title: Agent
sidebar_position: 2
sidebar_label: Agent
---

## FluxMeter

### Metrics

| Name       | Type      | Labels                                                                                             | Unit | Description |
| ---------- | --------- | -------------------------------------------------------------------------------------------------- | ---- | ----------- |
| flux_meter | Histogram | flux_meter_name, decision_type, response_status, http_status_code, feature_status, attribute_found | ms   | FluxMeter   |

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

| Name | Type | Labels | Unit | Description |
| ---- | ---- | ------ | ---- | ----------- |

## RateLimiter

### Metrics

| Name                 | Type    | Labels                                    | Unit | Description                                                        |
| -------------------- | ------- | ----------------------------------------- | ---- | ------------------------------------------------------------------ |
| rate_limiter_counter | Counter | policy_name, policy_hash, component_index |      | A counter measuring the number of times Rate Limiter was triggered |

### Labels

| Name            | Example | Description |
| --------------- | ------- | ----------- |
| policy_name     |         |             |
| policy_hash     |         |             |
| component_index |         |             |

## Classifier

### Metrics

| Name               | Type    | Labels                                     | Unit | Description                                                      |
| ------------------ | ------- | ------------------------------------------ | ---- | ---------------------------------------------------------------- |
| classifier_counter | Counter | policy_name, policy_hash, classifier_index |      | A counter measuring the number of times classifier was triggered |

### Labels

| Name             | Example | Description |
| ---------------- | ------- | ----------- |
| policy_name      |         |             |
| policy_hash      |         |             |
| classifier_index |         |             |

## Dataplane Summary

| Name | Type | Labels | Unit | Description |
| ---- | ---- | ------ | ---- | ----------- |
