local patch =
  {
    // Rename policylanguagev1FluxMeter to just FluxMeter
    local fluxMeter = super.policylanguagev1FluxMeter,
    policylanguagev1FluxMeter:: null,
    FluxMeter: fluxMeter {
      new(selector, attribute_key, buckets)::
        super.new()
        + fluxMeter.withSelector(selector)
        + fluxMeter.withAttributeKey(attribute_key)
        + fluxMeter.withStaticBuckets(buckets),
    },
  };

{
  v1+: patch,
}
