# Adapter Connector

| Status                   |                                                           |
| ------------------------ | --------------------------------------------------------- |
| Stability                | [Development]                                             |
| Supported pipeline types | See [Supported Pipeline Types](#supported-pipeline-types) |
| Distributions            | []                                                        |

The `adapter` connector converts between different OTEL signals.

## Supported Pipeline Types

| [Exporter Pipeline Type] | [Receiver Pipeline Type] |
| ------------------------ | ------------------------ |
| traces                   | logs                     |

## Configuration

If you are not already familiar with connectors, you may find it helpful to
first visit the [Connectors README].

The `adapter` connector does not have any configuration settings.

```yaml
receivers:
  foo:
exporters:
  bar:
connectors:
  adapter:
```

### Example Usage

Direct signals flow from metrics to logs and export.

```yaml
receivers:
  foo:
processors:
  attributes:
  batch:
exporters:
  bar:
connectors:
  adapter:
service:
  pipelines:
    metrics:
      receivers: [foo]
      processors: [attributes]
      exporters: [adapter]
    logs:
      receivers: [adapter]
      processors: [batch]
      exporters: [bar]
```

[Development]:
  https://github.com/open-telemetry/opentelemetry-collector#development
[Connectors README]:
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector/README.md
[Exporter Pipeline Type]:
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector/README.md#exporter-pipeline-type
[Receiver Pipeline Type]:
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector/README.md#receiver-pipeline-type
