{
  new():: {
  },
  withDeniedResponseStatusCode(denied_response_status_code):: {
    denied_response_status_code: denied_response_status_code,
  },
  withDeniedResponseStatusCodeMixin(denied_response_status_code):: {
    denied_response_status_code+: denied_response_status_code,
  },
  withRampMode(ramp_mode):: {
    ramp_mode: ramp_mode,
  },
  withRampModeMixin(ramp_mode):: {
    ramp_mode+: ramp_mode,
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
  withSessionLabelKey(session_label_key):: {
    session_label_key: session_label_key,
  },
  withSessionLabelKeyMixin(session_label_key):: {
    session_label_key+: session_label_key,
  },
}
