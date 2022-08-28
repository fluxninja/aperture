local patch =
  {
    Constant+: {
      local constant = self,
      new(output_port, value)::
        super.new()
        + constant.withValue(value)
        + constant.withOutPorts({ output: output_port }),
    },
  };

{
  v1+: patch,
}
