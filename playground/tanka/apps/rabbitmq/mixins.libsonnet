local rabbitmqApp = import 'apps/rabbitmq/main.libsonnet';

local rabbitmqAppMixin =
  rabbitmqApp {
    values+: {
      auth: {
        username: 'admin',
        password: 'secretpassword',
        erlangCookie: 'secretcookie',
      },
      persistence: {
        enabled: false,
      },
    },
  };

rabbitmqAppMixin
