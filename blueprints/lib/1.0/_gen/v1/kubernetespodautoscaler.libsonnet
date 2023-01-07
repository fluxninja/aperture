{
  new():: {
  },
  withKubernetesSelector(kubernetes_selector):: {
    kubernetes_selector: kubernetes_selector,
  },
  withKubernetesSelectorMixin(kubernetes_selector):: {
    kubernetes_selector+: kubernetes_selector,
  },
}
