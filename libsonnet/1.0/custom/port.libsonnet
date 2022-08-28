local patch =
  {
    Port+: {
      local port = self,
      new(signal_name)::
        super.new()
        + port.withSignalName(signal_name),
    },
  };

{
  v1+: patch,
}
