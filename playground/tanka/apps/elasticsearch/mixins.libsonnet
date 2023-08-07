local elasticsearchApp = import 'apps/elasticsearch/main.libsonnet';

local elasticsearchAppMixin =
  elasticsearchApp {
    values+: {
      extraConfig+: {
        network+: {
          host: '0.0.0.0',
        },
        discovery+: {
          seed_hosts: '',
        },
        xpack+: {
          monitoring+: {
            collection+: {
              enabled: true,
            },
          },
        },
      },
      master+: {
        extraRoles+: 'remote_cluster_client,ml',
        replicaCount: 1,
        masterOnly: false,
      },
      data+: {
        // extraRoles+: 'remote_cluster_client,ml',
        replicaCount: 0,
      },
      ingest+: {
        // extraRoles+: 'remote_cluster_client,ml',
        replicaCount: 0,
      },
      coordinating+: {
        // extraRoles+: 'remote_cluster_client,ml',
        replicaCount: 0,
      },
      security+: {
        elasticPassword: 'secretpassword',
      },
    },
  };

elasticsearchAppMixin
