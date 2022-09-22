{
  /**
   * Creates a [time series visualization][https://grafana.com/docs/grafana/latest/visualizations/time-series/]
   *
   */
  new(
    title,
    datasource=null,
    span=null,
    min_span=null,
    interval=null,
    axis_placement='auto',
    axis_label='',
    bar_alignment=0,
    color_mode='palette-classic',
    draw_style='line',
    fill_opacity=0,
    gradient_mode='none',
  ):: {
    _nextTarget: 0,

    title: title,
    type: 'timeseries',
    datasource: datasource,
    [if min_span != null then 'minSpan']: min_span,
    [if span != null then 'span']: span,
  },
  withTarget(target):: {
    local nextTarget = super._nextTarget,
    _nextTarget: nextTarget + 1,
    targets+: [target { refId: std.char(std.codepoint('A') + nextTarget) }],
  },
  withTargets(targets):: std.foldl(function(p, t) p.withTarget(t), targets, self),
  options: {
    withLegend(calcs=[], display_mode='list', placement='bottom'):: {
      legend: {
        calcs: calcs,
        displayMode: display_mode,
        placement: placement,
      },
    },
    withTooltip(mode='single', sort='none'):: {
      tooltip: {
        mode: mode,
        sort: sort,
      },
    },
  },
  defaults: {
    withColorMode(color_mode='palette-classic'):: {
      color: {
        mode: color_mode,
      },
    },
    withCustom(custom):: {
      custom: custom,
    },
    withCustomMixin(custom):: {
      custom+: custom,
    },
    withThresholds(thresholds):: {
      thresholds: thresholds,
    },
  },
  withDefaults(defaults):: {
    defaults: defaults,
  },
  withDefaultsMixin(defaults):: {
    defaults+: defaults,
  },
  fieldConfig: {
    withDefaults(defaults):: {
      defaults: defaults,
    },
    withDefaultsMixin(defaults):: {
      defaults+: defaults,
    },
    defaults: {
      withColor(color):: {
        color: color,
      },
      withCustom(custom):: {
        custom: custom,
      },
      withCustomMixin(custom):: {
        custom+: custom,
      },
      withThresholds(thresholds):: {
        thresholds: thresholds,
      },
      withMappings(mappings):: {
        mappings: mappings,
      },
      withUnit(unit):: {
        [if unit != '' then 'unit']: unit,
      },
    },
    withOverrides(overrides):: {
      overrides: overrides,
    },
    withOverridesMixin(overrides):: {
      overrides+: overrides,
    },
    overrides: {
      withMatcher(matcher):: {
        matcher: matcher,
      },
      withProperties(properties):: {
        properties: properties,
      },
      withPropertiesMixin(properties):: {
        properties+: properties,
      },
      withPropertiesOverride(properties):: {
        properties: properties,
      },
      withPropertiesOverrideMixin(properties):: {
        properties+: properties,
      },
    },
  },
  withFieldConfig(field_config):: {
    fieldConfig: field_config,
  },
  withOptions(options):: {
    options: options,
  },
  withOptionsMixin(options):: {
    options+: options,
  },
}
