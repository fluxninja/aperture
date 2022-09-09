{
  new():: {
  },
  withAttributeKey(attribute_key):: {
    attribute_key: attribute_key,
  },
  withAttributeKeyMixin(attribute_key):: {
    attribute_key+: attribute_key,
  },
  withHistogramBuckets(histogram_buckets):: {
    histogram_buckets:
      if std.isArray(histogram_buckets)
      then histogram_buckets
      else [histogram_buckets],
  },
  withHistogramBucketsMixin(histogram_buckets):: {
    histogram_buckets+: histogram_buckets,
  },
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
}
