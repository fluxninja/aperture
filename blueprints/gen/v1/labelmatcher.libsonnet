{
  new():: {
  },
  withExpression(expression):: {
    expression: expression,
  },
  withExpressionMixin(expression):: {
    expression+: expression,
  },
  withMatchLabels(match_labels):: {
    match_labels: match_labels,
  },
  withMatchLabelsMixin(match_labels):: {
    match_labels+: match_labels,
  },
}
