local patch = {
  local rateLimiter = super.languagev1RateLimiter,
  languagev1RateLimiter:: null,
  RateLimiter: rateLimiter {
    new(selector, limit_port, reset_interval, label_key)::
      rateLimiter.new()
      + rateLimiter.withSelector(selector)
      + rateLimiter.withLimitResetInterval(reset_interval)
      + rateLimiter.withLabelKey(label_key)
      + rateLimiter.withInPorts({ limit: limit_port }),
  },

};

{
  v1+: patch,
}
