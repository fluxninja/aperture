local patch =
  {
    SchedulerWorkloadAndLabelMatcher+: {
      local s = self,
      new(workload, label_matcher)::
        super.new()
        + s.withWorkload(workload)
        + s.withLabelMatcher(label_matcher),
    },
  };

{
  v1+: patch,
}
