{
  new():: {
  },
  withExtractor(extractor):: {
    extractor: extractor,
  },
  withExtractorMixin(extractor):: {
    extractor+: extractor,
  },
  withRego(rego):: {
    rego: rego,
  },
  withRegoMixin(rego):: {
    rego+: rego,
  },
  withTelemetry(telemetry):: {
    telemetry: telemetry,
  },
  withTelemetryMixin(telemetry):: {
    telemetry+: telemetry,
  },
}
