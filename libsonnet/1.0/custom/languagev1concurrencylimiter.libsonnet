local patch =
  {
    local r = super.languagev1ConcurrencyLimiter,
    languagev1ConcurrencyLimiter:: null,
    ConcurrencyLimiter: r,
  };

{
  v1+: patch,
}
