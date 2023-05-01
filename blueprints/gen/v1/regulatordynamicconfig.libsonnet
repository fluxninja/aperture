{
  new():: {
  },
  withEnableLabelValues(enable_label_values):: {
    enable_label_values:
      if std.isArray(enable_label_values)
      then enable_label_values
      else [enable_label_values],
  },
  withEnableLabelValuesMixin(enable_label_values):: {
    enable_label_values+: enable_label_values,
  },
}
