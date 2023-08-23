function(config, agent_group='default') {
  local jmx_regex = '(.*?:' + config.policy.jmx.jmx_metrics_port + ')',
  local app_address_regex = '(.*?:' + config.policy.jmx.app_server_port + ')',

  jmx_inframeter: {
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
                    names: [config.policy.jmx.app_namespace],
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
                  regex: jmx_regex,
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
                        names: [config.policy.jmx.app_namespace],
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
                    regex: app_address_regex,
                  },
                ],
            },
          ],
        },
      },
    },
  },
}
