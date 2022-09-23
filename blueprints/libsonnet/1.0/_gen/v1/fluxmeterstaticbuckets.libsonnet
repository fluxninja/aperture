{
  new():: {
  },
  withBuckets(buckets):: {
    buckets:
      if std.isArray(buckets)
      then buckets
      else [buckets],
  },
  withBucketsMixin(buckets):: {
    buckets+: buckets,
  },
}
