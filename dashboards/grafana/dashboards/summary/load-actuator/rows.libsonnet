local statPanel = import '../../../panels/stat.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';
local schedulerRowsFn = import '../scheduler/rows-fn.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(datasourceName, policyName, component, extraFilters={})
  local componentID = component.component.load_scheduler_component_id;
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID });

  local row1 = [
    statPanel(
      'Average Load Multiplier',
      datasourceName,
      query='avg(token_bucket_lm_ratio{%(filters)s})' % { filters: stringFilters },
      x=0,
      h=6,
      w=6
    ),
    statPanel(
      'Token Bucket Capacity',
      datasourceName,
      query='avg(token_bucket_capacity_total{%(filters)s})' % { filters: stringFilters },
      x=6,
      h=6,
      w=6
    ),
    statPanel(
      'Token Bucket FillRate',
      datasourceName,
      query='avg(token_bucket_fill_rate{%(filters)s})' % { filters: stringFilters },
      x=12,
      h=6,
      w=6
    ),
    statPanel(
      'Token Bucket Available Tokens',
      datasourceName,
      query='avg(token_bucket_available_tokens_total{%(filters)s})' % { filters: stringFilters },
      x=18,
      h=6,
      w=6
    ),
  ];

  schedulerRowsFn(datasourceName, policyName, componentID, extraFilters) + [row1]
