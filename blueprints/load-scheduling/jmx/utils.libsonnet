function(agent_group='default') {
  prometheus: {
    agent_group: agent_group,
    per_agent_group: true,
    pipeline: {
      receivers: [
        'prometheus',
      ],
    },
    receivers: {
      prometheus: {
        config: {
          scrape_configs: [
            {
              job_name: 'java-demo-app-jmx',
              scrape_interval: '10s',
              kubernetes_sd_configs: [
                {
                  role: 'pod',
                  namespaces: {
                    names: ['demoapp'],
                  },
                },
              ],
              relabel_configs: [
                {
                  source_labels: ['__meta_kubernetes_pod_annotation_prometheus_io_scrape'],
                  action: 'keep',
                  regex: true,
                },
                {
                  source_labels: ['__address__'],
                  action: 'keep',
                  regex: '(.*?:8087)',
                },
              ],
            },
            {
              job_name: 'java-demo-app-micrometer',
              scrape_interval: '10s',
              metrics_path: '/actuator/prometheus',
              kubernetes_sd_configs:
                [
                  {
                    role: 'pod',
                    namespaces:
                      {
                        names: ['demoapp'],
                      },
                  },
                ],
              relabel_configs:
                [
                  {
                    source_labels: ['__meta_kubernetes_pod_annotation_prometheus_io_scrape'],
                    action: 'keep',
                    regex: true,
                  },
                  {
                    source_labels: ['__address__'],
                    action: 'keep',
                    regex: '(.*?:8099)',
                  },
                ],
            },
          ],
        },
      },
    },
  },
}
