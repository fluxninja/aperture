{
  new():: {
  },
  withExtractor(extractor):: {
    extractor: extractor,
  },
  withExtractorMixin(extractor):: {
    extractor+: extractor,
  },
  withHidden(hidden):: {
    hidden: hidden,
  },
  withHiddenMixin(hidden):: {
    hidden+: hidden,
  },
  withPropagate(propagate):: {
    propagate: propagate,
  },
  withPropagateMixin(propagate):: {
    propagate+: propagate,
  },
  withRego(rego):: {
    rego: rego,
  },
  withRegoMixin(rego):: {
    rego+: rego,
  },
}
