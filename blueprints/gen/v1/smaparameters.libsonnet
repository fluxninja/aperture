{
  new():: {
  },
  withValidDuringWarmup(valid_during_warmup):: {
    valid_during_warmup: valid_during_warmup,
  },
  withValidDuringWarmupMixin(valid_during_warmup):: {
    valid_during_warmup+: valid_during_warmup,
  },
  withWindow(window):: {
    window: window,
  },
  withWindowMixin(window):: {
    window+: window,
  },
}
