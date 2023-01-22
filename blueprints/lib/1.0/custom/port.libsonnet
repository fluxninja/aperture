local patch =
  {
    Port: {
      withConstantSignal(constant_value):: {
        constant_value: constant_value,
      },
      withSignalName(signal_name):: {
        signal_name: signal_name,
      },
      withSignalNameMixin(signal_name):: {
        signal_name+: signal_name,
      },
    },
  };

{
  v1+: patch,
}
