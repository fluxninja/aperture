{
  new():: {
  },
  withDeniedResponseStatusCode(denied_response_status_code):: {
    denied_response_status_code: denied_response_status_code,
  },
  withDeniedResponseStatusCodeMixin(denied_response_status_code):: {
    denied_response_status_code+: denied_response_status_code,
  },
  withLabelKey(label_key):: {
    label_key: label_key,
  },
  withLabelKeyMixin(label_key):: {
    label_key+: label_key,
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
}
