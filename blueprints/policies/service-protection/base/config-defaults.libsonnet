local autoScalingDefaults = import '../../auto-scaling/base/config-defaults.libsonnet';

local selectors_defaults = [{
  service: '__REQUIRED_FIELD__',
  control_point: '__REQUIRED_FIELD__',
}];

local service_protection_core_defaults = {
  overload_confirmations: [],

  adaptive_load_scheduler: {
    load_scheduler: {
      selectors: selectors_defaults,
    },
    gradient: {
      slope: -1,
      min_gradient: 0.1,
      max_gradient: 1.0,
    },
    alerter: {
      alert_name: 'Load Throttling Event',
    },
    max_load_multiplier: 2.0,
    load_multiplier_linear_increment: 0.0025,
  },

  dry_run: false,
};

local auto_scaling_defaults = {
  dry_run: false,
  scaling_backend: {},
  promql_scale_out_controllers: [],
  promql_scale_in_controllers: [],
  scaling_parameters: autoScalingDefaults.policy.scaling_parameters,
  periodic_decrease: {
    period: '60s',
    scale_in_percentage: 10,
  },
};

local kubernetes_replicas_defaults = {
  kubernetes_object_selector: '__REQUIRED_FIELD__',
  min_replicas: '__REQUIRED_FIELD__',
  max_replicas: '__REQUIRED_FIELD__',
};

local kubeletstats_infra_meter = function(agent_group) {
  kubeletstats: {
    agent_group: agent_group,
    per_agent_group: true,
    pipeline: {
      processors: [
        'k8sattributes',
      ],
      receivers: [
        'kubeletstats',
      ],
    },
    processors: {
      k8sattributes: {
        auth_type: 'serviceAccount',
        passthrough: false,
        filter: {
          node_from_env_var: 'NODE_NAME',
        },
        extract: {
          metadata: [
            'k8s.cronjob.name',
            'k8s.daemonset.name',
            'k8s.deployment.name',
            'k8s.job.name',
            'k8s.namespace.name',
            'k8s.node.name',
            'k8s.pod.name',
            'k8s.pod.uid',
            'k8s.replicaset.name',
            'k8s.statefulset.name',
            'k8s.container.name',
          ],
        },
        pod_association: [
          {
            sources: [
              {
                from: 'resource_attribute',
                name: 'k8s.pod.uid',
              },
            ],
          },
        ],
      },
    },
    receivers: {
      kubeletstats: {
        collection_interval: '15s',
        auth_type: 'serviceAccount',
        endpoint: 'https://${NODE_NAME}:10250',
        insecure_skip_verify: true,
        metric_groups: [
          'pod',
          'container',
        ],
      },
    },
  },
};

{
  policy: {
    policy_name: '__REQUIRED_FIELD__',
    components: [],
    resources: {
      flow_control: {
        classifiers: [],
      },
    },
    evaluation_interval: '1s',
    service_protection_core: service_protection_core_defaults,
    auto_scaling: auto_scaling_defaults {
      scaling_backend+: {
        kubernetes_replicas: kubernetes_replicas_defaults,
      },
    },
  },

  dashboard: {
    refresh_interval: '5s',
    time_from: 'now-15m',
    time_to: 'now',
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
    variant_name: 'Service Protection',
  },

  selectors: selectors_defaults,

  auto_scaling_pods: auto_scaling_defaults {
    scaling_backend+: {
      kubernetes_replicas: kubernetes_replicas_defaults,
    },
  },

  kubeletstats_infra_meter: kubeletstats_infra_meter,
}
