{
  new():: {
  },
  withLabelKey(label_key):: {
    label_key: label_key,
  },
  withLabelKeyMixin(label_key):: {
    label_key+: label_key,
  },
  withLazySync(lazy_sync):: {
    lazy_sync: lazy_sync,
  },
  withLazySyncMixin(lazy_sync):: {
    lazy_sync+: lazy_sync,
  },
  withLimitResetInterval(limit_reset_interval):: {
    limit_reset_interval: limit_reset_interval,
  },
  withLimitResetIntervalMixin(limit_reset_interval):: {
    limit_reset_interval+: limit_reset_interval,
  },
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
  withTokensLabelKey(tokens_label_key):: {
    tokens_label_key: tokens_label_key,
  },
  withTokensLabelKeyMixin(tokens_label_key):: {
    tokens_label_key+: tokens_label_key,
  },
}
