{
  new():: {
  },
  withDryRun(dry_run):: {
    dry_run: dry_run,
  },
  withDryRunMixin(dry_run):: {
    dry_run+: dry_run,
  },
  withDryRunConfigKey(dry_run_config_key):: {
    dry_run_config_key: dry_run_config_key,
  },
  withDryRunConfigKeyMixin(dry_run_config_key):: {
    dry_run_config_key+: dry_run_config_key,
  },
  withKubernetesObjectSelector(kubernetes_object_selector):: {
    kubernetes_object_selector: kubernetes_object_selector,
  },
  withKubernetesObjectSelectorMixin(kubernetes_object_selector):: {
    kubernetes_object_selector+: kubernetes_object_selector,
  },
}
