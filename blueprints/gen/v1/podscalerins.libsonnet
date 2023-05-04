{
  new():: {
  },
  withReplicas(replicas):: {
    replicas: replicas,
  },
  withReplicasMixin(replicas):: {
    replicas+: replicas,
  },
}
