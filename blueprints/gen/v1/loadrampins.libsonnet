{
  new():: {
  },
  withBackward(backward):: {
    backward: backward,
  },
  withBackwardMixin(backward):: {
    backward+: backward,
  },
  withForward(forward):: {
    forward: forward,
  },
  withForwardMixin(forward):: {
    forward+: forward,
  },
  withReset(reset):: {
    reset: reset,
  },
  withResetMixin(reset):: {
    reset+: reset,
  },
}
