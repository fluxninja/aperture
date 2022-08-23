local k = import 'k.libsonnet';

local bookinfoApp = import 'apps/bookinfo/main.libsonnet';

local bookinfoMixin =
  bookinfoApp {
    values+: {
    },
  };

bookinfoMixin
