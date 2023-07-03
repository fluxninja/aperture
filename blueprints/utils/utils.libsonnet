{
  resources: function(resources) {
    updatedResources:
      if std.objectHas(resources, 'telemetry_collectors') then resources {
        telemetry_collectors:: {},
        local infraMeters = if std.objectHas(resources, 'infra_meters') then resources.infra_meters else {},
        local addAgentGroup(collector, infraMeters) = {
          [k]: if std.objectHas(infraMeters, k) then infraMeters[k] { agent_group: collector.agent_group } else error 'Invalid key'
          for k in std.objectFields(infraMeters)
        },
        local merge(infraMetersInitial, infraMeters) = infraMetersInitial + {
          [k]: if std.objectHas(infraMetersInitial, k) then
            error 'Conflicting keys found in policy.resources.infra_meters and policy.resources.telemetry_collectors'
          else infraMeters[k]
          for k in std.objectFields(infraMeters)
        },
        infra_meters: std.foldl(merge, [addAgentGroup(collector, collector.infra_meters) for collector in resources.telemetry_collectors], infraMeters),
      } else resources,
  },
}
