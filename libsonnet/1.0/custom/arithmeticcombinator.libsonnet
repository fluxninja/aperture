local patch =
  {
    ArithmeticCombinator+: {
      local combinator = self,
      local op_(operator, lhs, rhs, output) =
        combinator.new()
        + combinator.withOperator(operator)
        + combinator.withInPorts({ lhs: lhs, rhs: rhs })
        + combinator.withOutPorts({ output: output }),
      mul(lhs, rhs, output):: op_('mul', lhs, rhs, output),
      div(lhs, rhs, output):: op_('div', lhs, rhs, output),
      add(lhs, rhs, output):: op_('add', lhs, rhs, output),
      sub(lhs, rhs, output):: op_('sub', lhs, rhs, output),
    },
  };

{
  v1+: patch,
}
