{
  new():: {
  },
  withQuery(query):: {
    query: query,
  },
  withQueryMixin(query):: {
    query+: query,
  },
  withSource(source):: {
    source: source,
  },
  withSourceMixin(source):: {
    source+: source,
  },
}
