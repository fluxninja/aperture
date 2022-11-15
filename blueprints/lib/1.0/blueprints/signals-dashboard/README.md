# Signals Dashboard

## Configuration

[configuration]: # Configuration Marker

### Common

| Parameter Name      | Parameter Type | Default      | Description         |
| ------------------- | -------------- | ------------ | ------------------- |
| `common.policyName` | `string`       | `(required)` | Name of the policy. |

### Dashboard

| Parameter Name              | Parameter Type | Default     | Description                            |
| --------------------------- | -------------- | ----------- | -------------------------------------- |
| `dashboard.refreshInterval` | `string`       | `"10s"`     | Refresh interval for dashboard panels. |
| `dashboard.timeFrom`        | `string`       | `"now-30m"` | From time of dashboard.                |
| `dashboard.timeTo`          | `string`       | `"now"`     | To time of dashboard.                  |

#### Datasource

| Parameter Name                     | Parameter Type | Default         | Description              |
| ---------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`        | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filterRegex` | `string`       | `""`            | Datasource filter regex. |
