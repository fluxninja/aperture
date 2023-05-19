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
}
