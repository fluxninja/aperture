{
  new():: {
  },
  withDesiredReplicas(desired_replicas):: {
    desired_replicas: desired_replicas,
  },
  withDesiredReplicasMixin(desired_replicas):: {
    desired_replicas+: desired_replicas,
  },
}
