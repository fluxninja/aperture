{
  new():: {
  },
  withTelemetry(telemetry):: {
    telemetry: telemetry,
  },
  withTelemetryMixin(telemetry):: {
    telemetry+: telemetry,
  },
}
