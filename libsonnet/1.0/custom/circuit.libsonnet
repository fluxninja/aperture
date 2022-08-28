local patch =
  {
    Circuit+: {
      local circuit = self,
      new(evaluation_interval)::
        super.new()
        + circuit.withEvaluationInterval(evaluation_interval),
    },
  };

{
  v1+: patch,
}
