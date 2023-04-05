{
  new():: {
  },
  withActualReplicas(actual_replicas):: {
    actual_replicas: actual_replicas,
  },
  withActualReplicasMixin(actual_replicas):: {
    actual_replicas+: actual_replicas,
  },
  withConfiguredReplicas(configured_replicas):: {
    configured_replicas: configured_replicas,
  },
  withConfiguredReplicasMixin(configured_replicas):: {
    configured_replicas+: configured_replicas,
  },
  withDesiredReplicas(desired_replicas):: {
    desired_replicas: desired_replicas,
  },
  withDesiredReplicasMixin(desired_replicas):: {
    desired_replicas+: desired_replicas,
  },
}
