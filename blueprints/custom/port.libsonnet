local is_nan(val) =
  if val == 'NaN' || val == '-Inf' || val == '+Inf' then true else false;

local patch =
  {
    Port: {
      withConstantSignal(constant_signal):: {
        constant_signal: {
          [if is_nan(constant_signal) then 'special_value' else 'value']: constant_signal,
        },
      },
      withConstantSignalMixin(constant_signal):: {
        constant_signal: {
          [if is_nan(constant_signal) then 'special_value' else 'value']+: constant_signal,
        },
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
