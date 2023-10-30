local serviceProtectionDefaults = import '../common-overload-protection/config-defaults.libsonnet';

serviceProtectionDefaults {

  dashboard+: {
    title: 'Aperture Service Protection for CPU Overload',
    variant_name: 'CPU Overload Detection',
  },
}
