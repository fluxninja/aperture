{
  new():: {
  },
  withKubernetesReplicas(kubernetes_replicas):: {
    kubernetes_replicas: kubernetes_replicas,
  },
  withKubernetesReplicasMixin(kubernetes_replicas):: {
    kubernetes_replicas+: kubernetes_replicas,
  },
}
