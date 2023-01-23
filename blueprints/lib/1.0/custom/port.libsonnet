local patch =
  {
    Port: {
      withConstantSignal(constant_signal):: {
        constant_signal: constant_signal,
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
