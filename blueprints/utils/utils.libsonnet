local kubeletstats_infra_meter(agent_group, selectors) = {
  kubeletstats: {
    agent_group: agent_group,
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
        selectors: if std.isArray(selectors) then selectors else [selectors],
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

local merge(infraMetersInitial, infraMeters) = infraMetersInitial + {
  [k]: if std.objectHas(infraMetersInitial, k) then
    error std.format('Conflicting key "%s" found for policy.resources.infra_meters', [k])
  else infraMeters[k]
  for k in std.objectFields(infraMeters)
};

local add_kubeletstats_infra_meter(infra_meters, agent_group='default', selector={}) =
  merge(
    infra_meters,
    kubeletstats_infra_meter(agent_group, selector)
  );

{
  add_kubeletstats_infra_meter: add_kubeletstats_infra_meter,
}
