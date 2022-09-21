local patch =
  {
    // Rename policylanguagev1FluxMeter to just FluxMeter
    local fluxMeter = super.policylanguagev1FluxMeter,
    policylanguagev1FluxMeter:: null,
    FluxMeter: fluxMeter {
      new(selector)::
        super.new()
        + fluxMeter.withSelector(selector),
    },
  };

{
  v1+: patch,
}
