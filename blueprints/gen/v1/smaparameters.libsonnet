{
  new():: {
  },
  withSmaWindow(sma_window):: {
    sma_window: sma_window,
  },
  withSmaWindowMixin(sma_window):: {
    sma_window+: sma_window,
  },
  withValidDuringWarmup(valid_during_warmup):: {
    valid_during_warmup: valid_during_warmup,
  },
  withValidDuringWarmupMixin(valid_during_warmup):: {
    valid_during_warmup+: valid_during_warmup,
  },
}
