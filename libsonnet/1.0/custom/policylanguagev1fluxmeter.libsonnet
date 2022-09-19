local patch =
  {
    // Rename policylanguagev1FluxMeter to just FluxMeter
    local fluxMeter = super.policylanguagev1FluxMeter,
    policylanguagev1FluxMeter:: null,
    FluxMeter: fluxMeter {
      new(selector, buckets)::
        super.new()
        + fluxMeter.withSelector(selector)
        + fluxMeter.withBuckets(buckets),
    },
  };

{
  v1+: patch,
}
