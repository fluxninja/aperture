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
