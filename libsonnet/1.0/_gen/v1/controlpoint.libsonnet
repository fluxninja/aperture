{
  new():: {
  },
  withFeature(feature):: {
    feature: feature,
  },
  withFeatureMixin(feature):: {
    feature+: feature,
  },
  withTraffic(traffic):: {
    traffic: traffic,
  },
  withTrafficMixin(traffic):: {
    traffic+: traffic,
  },
}
