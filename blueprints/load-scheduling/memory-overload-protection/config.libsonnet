local serviceProtectionDefaults = import '../common-overload-protection/config-defaults.libsonnet';

serviceProtectionDefaults {
  dashboard+: {
    title: 'Aperture Service Protection for Memory Overload',
    variant_name: 'Memory Overload Detection',
  },
}
