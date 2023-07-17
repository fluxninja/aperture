local elasticsearchApp = import 'apps/elasticsearch/main.libsonnet';

local elasticsearchAppMixin =
  elasticsearchApp {
    values+: {
      extraEnvVars: [
        {
          name: 'ELASTICSEARCH_USERNAME',
          value: 'elasticsearch',
        },
        {
          name: 'ELASTICSEARCH_PASSWORD',
          value: 'ThisIsASuperSecurePassword!',
        },
      ],
      master+: {
        extraRoles+: 'remote_cluster_client,ml',
      },
      data+: {
        extraRoles+: 'remote_cluster_client,ml',
      },
      ingest+: {
        extraRoles+: 'remote_cluster_client,ml',
      },
      coordinating+: {
        extraRoles+: 'remote_cluster_client,ml',
      },
    },
  };

elasticsearchAppMixin
