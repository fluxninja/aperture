{
  new():: {
  },
  withProcessors(processors):: {
    processors:
      if std.isArray(processors)
      then processors
      else [processors],
  },
  withProcessorsMixin(processors):: {
    processors+: processors,
  },
  withReceivers(receivers):: {
    receivers:
      if std.isArray(receivers)
      then receivers
      else [receivers],
  },
  withReceiversMixin(receivers):: {
    receivers+: receivers,
  },
}
