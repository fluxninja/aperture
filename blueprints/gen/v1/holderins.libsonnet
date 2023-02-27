{
  new():: {
  },
  withInput(input):: {
    input: input,
  },
  withInputMixin(input):: {
    input+: input,
  },
  withReset(reset):: {
    reset: reset,
  },
  withResetMixin(reset):: {
    reset+: reset,
  },
}
