{
  new():: {
  },
  withInPortsList(in_ports_list):: {
    in_ports_list:
      if std.isArray(in_ports_list)
      then in_ports_list
      else [in_ports_list],
  },
  withInPortsListMixin(in_ports_list):: {
    in_ports_list+: in_ports_list,
  },
}
