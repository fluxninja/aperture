{
  new():: {
  },
  withBucketCapacity(bucket_capacity):: {
    bucket_capacity: bucket_capacity,
  },
  withBucketCapacityMixin(bucket_capacity):: {
    bucket_capacity+: bucket_capacity,
  },
  withLeakAmount(leak_amount):: {
    leak_amount: leak_amount,
  },
  withLeakAmountMixin(leak_amount):: {
    leak_amount+: leak_amount,
  },
  withLeakIntervalMs(leak_interval_ms):: {
    leak_interval_ms: leak_interval_ms,
  },
  withLeakIntervalMsMixin(leak_interval_ms):: {
    leak_interval_ms+: leak_interval_ms,
  },
}
