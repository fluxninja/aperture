local is_special(val) =
  if std.isString(val) && (val == 'NaN' || val == '-Inf' || val == '+Inf') then true else false;

local patch =
  {
    Port: {
      withConstantSignal(constant_signal):: {
        constant_signal: {
          [if is_special(constant_signal) then 'special_value' else 'value']: constant_signal,
        },
      },
      withConstantSignalMixin(constant_signal):: {
        constant_signal: {
          [if is_special(constant_signal) then 'special_value' else 'value']+: constant_signal,
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
