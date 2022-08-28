{
  new():: {
  },
  withAll(all):: {
    all: all,
  },
  withAllMixin(all):: {
    all+: all,
  },
  withAny(any):: {
    any: any,
  },
  withAnyMixin(any):: {
    any+: any,
  },
  withLabelEquals(label_equals):: {
    label_equals: label_equals,
  },
  withLabelEqualsMixin(label_equals):: {
    label_equals+: label_equals,
  },
  withLabelExists(label_exists):: {
    label_exists: label_exists,
  },
  withLabelExistsMixin(label_exists):: {
    label_exists+: label_exists,
  },
  withLabelMatches(label_matches):: {
    label_matches: label_matches,
  },
  withLabelMatchesMixin(label_matches):: {
    label_matches+: label_matches,
  },
  withNot(not):: {
    not: not,
  },
  withNotMixin(not):: {
    not+: not,
  },
}
