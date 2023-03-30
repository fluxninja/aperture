{
  new():: {
  },
  withQuery(query):: {
    query: query,
  },
  withQueryMixin(query):: {
    query+: query,
  },
  withTelemetry(telemetry):: {
    telemetry: telemetry,
  },
  withTelemetryMixin(telemetry):: {
    telemetry+: telemetry,
  },
}
