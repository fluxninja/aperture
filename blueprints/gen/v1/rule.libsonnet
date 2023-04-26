{
  new():: {
  },
  withExtractor(extractor):: {
    extractor: extractor,
  },
  withExtractorMixin(extractor):: {
    extractor+: extractor,
  },
  withTelemetry(telemetry):: {
    telemetry: telemetry,
  },
  withTelemetryMixin(telemetry):: {
    telemetry+: telemetry,
  },
}
