{
  new():: {
  },
  withOutPortsList(out_ports_list):: {
    out_ports_list:
      if std.isArray(out_ports_list)
      then out_ports_list
      else [out_ports_list],
  },
  withOutPortsListMixin(out_ports_list):: {
    out_ports_list+: out_ports_list,
  },
}
