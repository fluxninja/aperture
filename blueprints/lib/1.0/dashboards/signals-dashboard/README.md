# Signals Dashboard

This blueprint provides a [policy monitoring](/reference/policies/monitoring.md)
dashboard that visualizes Signals flowing through the
[Circuit](/concepts/policy/circuit.md).

## Configuration

<!-- Configuration Marker -->

### Common

| Parameter Name       | Parameter Type | Default      | Description         |
| -------------------- | -------------- | ------------ | ------------------- |
| `common.policy_name` | `string`       | `(required)` | Name of the policy. |

### Dashboard

| Parameter Name               | Parameter Type | Default     | Description                            |
| ---------------------------- | -------------- | ----------- | -------------------------------------- |
| `dashboard.refresh_interval` | `string`       | `"10s"`     | Refresh interval for dashboard panels. |
| `dashboard.time_from`        | `string`       | `"now-30m"` | From time of dashboard.                |
| `dashboard.time_to`          | `string`       | `"now"`     | To time of dashboard.                  |

#### Datasource

| Parameter Name                      | Parameter Type | Default         | Description              |
| ----------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`         | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filter_regex` | `string`       | `""`            | Datasource filter regex. |
