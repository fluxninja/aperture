{
  new():: {
  },
  withPriority(priority):: {
    priority: priority,
  },
  withPriorityMixin(priority):: {
    priority+: priority,
  },
  withTokens(tokens):: {
    tokens: tokens,
  },
  withTokensMixin(tokens):: {
    tokens+: tokens,
  },
}
