{
  new():: {
  },
  withComponents(components):: {
    components:
      if std.isArray(components)
      then components
      else [components],
  },
  withComponentsMixin(components):: {
    components+: components,
  },
  withInPortsMap(in_ports_map):: {
    in_ports_map: in_ports_map,
  },
  withInPortsMapMixin(in_ports_map):: {
    in_ports_map+: in_ports_map,
  },
  withName(name):: {
    name: name,
  },
  withNameMixin(name):: {
    name+: name,
  },
  withOutPortsMap(out_ports_map):: {
    out_ports_map: out_ports_map,
  },
  withOutPortsMapMixin(out_ports_map):: {
    out_ports_map+: out_ports_map,
  },
}
