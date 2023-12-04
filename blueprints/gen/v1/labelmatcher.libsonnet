{
  new():: {
  },
  withExpression(expression):: {
    expression: expression,
  },
  withExpressionMixin(expression):: {
    expression+: expression,
  },
  withMatchExpressions(match_expressions):: {
    match_expressions:
      if std.isArray(match_expressions)
      then match_expressions
      else [match_expressions],
  },
  withMatchExpressionsMixin(match_expressions):: {
    match_expressions+: match_expressions,
  },
  withMatchLabels(match_labels):: {
    match_labels: match_labels,
  },
  withMatchLabelsMixin(match_labels):: {
    match_labels+: match_labels,
  },
  withMatchList(match_list):: {
    match_list:
      if std.isArray(match_list)
      then match_list
      else [match_list],
  },
  withMatchListMixin(match_list):: {
    match_list+: match_list,
  },
}
