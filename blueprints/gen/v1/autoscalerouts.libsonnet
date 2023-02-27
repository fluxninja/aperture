{
  new():: {
  },
  withActualScale(actual_scale):: {
    actual_scale: actual_scale,
  },
  withActualScaleMixin(actual_scale):: {
    actual_scale+: actual_scale,
  },
  withConfiguredScale(configured_scale):: {
    configured_scale: configured_scale,
  },
  withConfiguredScaleMixin(configured_scale):: {
    configured_scale+: configured_scale,
  },
  withDesiredScale(desired_scale):: {
    desired_scale: desired_scale,
  },
  withDesiredScaleMixin(desired_scale):: {
    desired_scale+: desired_scale,
  },
}
