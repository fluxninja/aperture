---
title: Controller
sidebar_position: 3
sidebar_label: Controller
---

## Signal

### Metrics

| Name           | Type    | Labels                   | Unit | Description               |
| -------------- | ------- | ------------------------ | ---- | ------------------------- |
| signal_reading | Summary | signal_name, policy_name |      | The reading from a signal |

### Labels

| Name        | Example                  | Description                                                                                    |
| ----------- | ------------------------ | ---------------------------------------------------------------------------------------------- |
| signal_name | LATENCY_EMA, IS_OVERLOAD | Name of the signal provided in policy.                                                         |
| policy_name | service1-demo-app        | Name of the policy.                                                                            |
| valid       | true, false              | Label for specifying if metric is valid. A metric may be invalid if signal reading is invalid. |
