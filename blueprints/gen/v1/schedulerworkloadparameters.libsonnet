{
  new():: {
  },
  withPriority(priority):: {
    priority: priority,
  },
  withPriorityMixin(priority):: {
    priority+: priority,
  },
  withQueueTimeout(queue_timeout):: {
    queue_timeout: queue_timeout,
  },
  withQueueTimeoutMixin(queue_timeout):: {
    queue_timeout+: queue_timeout,
  },
  withTokens(tokens):: {
    tokens: tokens,
  },
  withTokensMixin(tokens):: {
    tokens+: tokens,
  },
}
