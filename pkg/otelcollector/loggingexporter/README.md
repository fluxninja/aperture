# Logging Exporter

This is a fork of Logging Exporter from OTEL Collector with added support for
Aperture logging instead of zap.

Exports data to the console via Aperture logger.

Supported pipeline types: traces, metrics, logs

## Getting Started

The following settings are optional:

- `sampling_initial` (default = `2`): number of messages initially logged each
  second.
- `sampling_thereafter` (default = `500`): sampling rate after the initial
  messages are logged (every Mth message is logged). Refer to
  [Zap docs](https://godoc.org/go.uber.org/zap/zapcore#NewSampler) for more
  details. on how sampling parameters impact number of messages.

Example:

```yaml
exporters:
  aperturelogging:
    sampling_initial: 5
    sampling_thereafter: 200
```
