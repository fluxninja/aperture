{
  new():: {
  },
  withCircuit(circuit):: {
    circuit: circuit,
  },
  withCircuitMixin(circuit):: {
    circuit+: circuit,
  },
  withResources(resources):: {
    resources: resources,
  },
  withResourcesMixin(resources):: {
    resources+: resources,
  },
}
