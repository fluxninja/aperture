local postgresqlApp = import 'apps/postgresql/main.libsonnet';

local postgresqlAppMixin =
  postgresqlApp {
    values+: {
      auth: {
        username: 'postgres',
        password: 'secretpassword',
      },
    },
  };

postgresqlAppMixin
