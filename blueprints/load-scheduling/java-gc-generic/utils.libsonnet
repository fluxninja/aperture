function(config, agent_group='default') {
  // Define the target address directly
  local jmx_target_address = config.policy.jmx.jmx_host + ':' + config.policy.jmx.jmx_prometheus_port,

  jmx_inframeter: {
    agent_group: agent_group,
    per_agent_group: false,
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
              static_configs: [
                {
                  targets: [jmx_target_address],
                },
              ],
            },
          ],
        },
      },
    },
  },
}
