{
  new():: {
  },
  withKubernetesObjectSelector(kubernetes_object_selector):: {
    kubernetes_object_selector: kubernetes_object_selector,
  },
  withKubernetesObjectSelectorMixin(kubernetes_object_selector):: {
    kubernetes_object_selector+: kubernetes_object_selector,
  },
  withMaxReplicas(max_replicas):: {
    max_replicas: max_replicas,
  },
  withMaxReplicasMixin(max_replicas):: {
    max_replicas+: max_replicas,
  },
  withMinReplicas(min_replicas):: {
    min_replicas: min_replicas,
  },
  withMinReplicasMixin(min_replicas):: {
    min_replicas+: min_replicas,
  },
}
