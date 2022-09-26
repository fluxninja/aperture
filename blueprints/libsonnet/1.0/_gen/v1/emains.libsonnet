{
  new():: {
  },
  withInput(input):: {
    input: input,
  },
  withInputMixin(input):: {
    input+: input,
  },
  withMaxEnvelope(max_envelope):: {
    max_envelope: max_envelope,
  },
  withMaxEnvelopeMixin(max_envelope):: {
    max_envelope+: max_envelope,
  },
  withMinEnvelope(min_envelope):: {
    min_envelope: min_envelope,
  },
  withMinEnvelopeMixin(min_envelope):: {
    min_envelope+: min_envelope,
  },
}
