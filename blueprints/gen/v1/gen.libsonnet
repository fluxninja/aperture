{
  AdaptiveLoadScheduler: import 'adaptiveloadscheduler.libsonnet',
  AdaptiveLoadSchedulerIns: import 'adaptiveloadschedulerins.libsonnet',
  AdaptiveLoadSchedulerOuts: import 'adaptiveloadschedulerouts.libsonnet',
  AdaptiveLoadSchedulerParameters: import 'adaptiveloadschedulerparameters.libsonnet',
  AddressExtractor: import 'addressextractor.libsonnet',
  Alerter: import 'alerter.libsonnet',
  AlerterIns: import 'alerterins.libsonnet',
  AlerterParameters: import 'alerterparameters.libsonnet',
  And: import 'and.libsonnet',
  AndIns: import 'andins.libsonnet',
  AndOuts: import 'andouts.libsonnet',
  ArithmeticCombinator: import 'arithmeticcombinator.libsonnet',
  ArithmeticCombinatorIns: import 'arithmeticcombinatorins.libsonnet',
  ArithmeticCombinatorOuts: import 'arithmeticcombinatorouts.libsonnet',
  AutoScale: import 'autoscale.libsonnet',
  AutoScaler: import 'autoscaler.libsonnet',
  AutoScalerScalingBackend: import 'autoscalerscalingbackend.libsonnet',
  AutoScalerScalingBackendKubernetesReplicas: import 'autoscalerscalingbackendkubernetesreplicas.libsonnet',
  AutoScalerScalingBackendKubernetesReplicasOuts: import 'autoscalerscalingbackendkubernetesreplicasouts.libsonnet',
  AutoScalerScalingParameters: import 'autoscalerscalingparameters.libsonnet',
  Circuit: import 'circuit.libsonnet',
  Classifier: import 'classifier.libsonnet',
  Component: import 'component.libsonnet',
  ConstantSignal: import 'constantsignal.libsonnet',
  Decider: import 'decider.libsonnet',
  DeciderIns: import 'deciderins.libsonnet',
  DeciderOuts: import 'deciderouts.libsonnet',
  DecreasingGradient: import 'decreasinggradient.libsonnet',
  DecreasingGradientIns: import 'decreasinggradientins.libsonnet',
  DecreasingGradientParameters: import 'decreasinggradientparameters.libsonnet',
  Differentiator: import 'differentiator.libsonnet',
  DifferentiatorIns: import 'differentiatorins.libsonnet',
  DifferentiatorOuts: import 'differentiatorouts.libsonnet',
  EMA: import 'ema.libsonnet',
  EMAIns: import 'emains.libsonnet',
  EMAOuts: import 'emaouts.libsonnet',
  EMAParameters: import 'emaparameters.libsonnet',
  EqualsMatchExpression: import 'equalsmatchexpression.libsonnet',
  Extractor: import 'extractor.libsonnet',
  Extrapolator: import 'extrapolator.libsonnet',
  ExtrapolatorIns: import 'extrapolatorins.libsonnet',
  ExtrapolatorOuts: import 'extrapolatorouts.libsonnet',
  ExtrapolatorParameters: import 'extrapolatorparameters.libsonnet',
  FirstValid: import 'firstvalid.libsonnet',
  FirstValidIns: import 'firstvalidins.libsonnet',
  FirstValidOuts: import 'firstvalidouts.libsonnet',
  FlowControl: import 'flowcontrol.libsonnet',
  FlowControlResources: import 'flowcontrolresources.libsonnet',
  FluxMeter: import 'fluxmeter.libsonnet',
  FluxMeterExponentialBuckets: import 'fluxmeterexponentialbuckets.libsonnet',
  FluxMeterExponentialBucketsRange: import 'fluxmeterexponentialbucketsrange.libsonnet',
  FluxMeterLinearBuckets: import 'fluxmeterlinearbuckets.libsonnet',
  FluxMeterStaticBuckets: import 'fluxmeterstaticbuckets.libsonnet',
  GradientController: import 'gradientcontroller.libsonnet',
  GradientControllerDynamicConfig: import 'gradientcontrollerdynamicconfig.libsonnet',
  GradientControllerIns: import 'gradientcontrollerins.libsonnet',
  GradientControllerOuts: import 'gradientcontrollerouts.libsonnet',
  GradientControllerParameters: import 'gradientcontrollerparameters.libsonnet',
  Holder: import 'holder.libsonnet',
  HolderIns: import 'holderins.libsonnet',
  HolderOuts: import 'holderouts.libsonnet',
  InPort: import 'inport.libsonnet',
  IncreasingGradient: import 'increasinggradient.libsonnet',
  IncreasingGradientIns: import 'increasinggradientins.libsonnet',
  IncreasingGradientParameters: import 'increasinggradientparameters.libsonnet',
  InfraMeter: import 'inframeter.libsonnet',
  InfraMeterMetricsPipeline: import 'inframetermetricspipeline.libsonnet',
  Integrator: import 'integrator.libsonnet',
  IntegratorIns: import 'integratorins.libsonnet',
  IntegratorOuts: import 'integratorouts.libsonnet',
  Inverter: import 'inverter.libsonnet',
  InverterIns: import 'inverterins.libsonnet',
  InverterOuts: import 'inverterouts.libsonnet',
  JSONExtractor: import 'jsonextractor.libsonnet',
  JWTExtractor: import 'jwtextractor.libsonnet',
  K8sLabelMatcherRequirement: import 'k8slabelmatcherrequirement.libsonnet',
  KubernetesObjectSelector: import 'kubernetesobjectselector.libsonnet',
  LabelMatcher: import 'labelmatcher.libsonnet',
  LoadRamp: import 'loadramp.libsonnet',
  LoadRampIns: import 'loadrampins.libsonnet',
  LoadRampOuts: import 'loadrampouts.libsonnet',
  LoadRampParameters: import 'loadrampparameters.libsonnet',
  LoadRampParametersStep: import 'loadrampparametersstep.libsonnet',
  LoadRampSeries: import 'loadrampseries.libsonnet',
  LoadRampSeriesIns: import 'loadrampseriesins.libsonnet',
  LoadRampSeriesLoadRampInstance: import 'loadrampseriesloadrampinstance.libsonnet',
  LoadRampSeriesParameters: import 'loadrampseriesparameters.libsonnet',
  LoadScheduler: import 'loadscheduler.libsonnet',
  LoadSchedulerDynamicConfig: import 'loadschedulerdynamicconfig.libsonnet',
  LoadSchedulerIns: import 'loadschedulerins.libsonnet',
  LoadSchedulerOuts: import 'loadschedulerouts.libsonnet',
  LoadSchedulerParameters: import 'loadschedulerparameters.libsonnet',
  MatchExpression: import 'matchexpression.libsonnet',
  MatchExpressionList: import 'matchexpressionlist.libsonnet',
  MatchesMatchExpression: import 'matchesmatchexpression.libsonnet',
  Max: import 'max.libsonnet',
  MaxIns: import 'maxins.libsonnet',
  MaxOuts: import 'maxouts.libsonnet',
  Min: import 'min.libsonnet',
  MinIns: import 'minins.libsonnet',
  MinOuts: import 'minouts.libsonnet',
  NestedCircuit: import 'nestedcircuit.libsonnet',
  NestedSignalEgress: import 'nestedsignalegress.libsonnet',
  NestedSignalEgressIns: import 'nestedsignalegressins.libsonnet',
  NestedSignalIngress: import 'nestedsignalingress.libsonnet',
  NestedSignalIngressOuts: import 'nestedsignalingressouts.libsonnet',
  Or: import 'or.libsonnet',
  OrIns: import 'orins.libsonnet',
  OrOuts: import 'orouts.libsonnet',
  OutPort: import 'outport.libsonnet',
  PathTemplateMatcher: import 'pathtemplatematcher.libsonnet',
  PeriodicDecrease: import 'periodicdecrease.libsonnet',
  PeriodicDecreaseParameters: import 'periodicdecreaseparameters.libsonnet',
  PodScaler: import 'podscaler.libsonnet',
  PodScalerIns: import 'podscalerins.libsonnet',
  PodScalerOuts: import 'podscalerouts.libsonnet',
  Policy: import 'policy.libsonnet',
  PromQL: import 'promql.libsonnet',
  PromQLOuts: import 'promqlouts.libsonnet',
  PulseGenerator: import 'pulsegenerator.libsonnet',
  PulseGeneratorOuts: import 'pulsegeneratorouts.libsonnet',
  Query: import 'query.libsonnet',
  RateLimiter: import 'ratelimiter.libsonnet',
  RateLimiterDynamicConfig: import 'ratelimiterdynamicconfig.libsonnet',
  RateLimiterIns: import 'ratelimiterins.libsonnet',
  RateLimiterOverride: import 'ratelimiteroverride.libsonnet',
  RateLimiterParameters: import 'ratelimiterparameters.libsonnet',
  RateLimiterParametersLazySync: import 'ratelimiterparameterslazysync.libsonnet',
  Rego: import 'rego.libsonnet',
  RegoLabelProperties: import 'regolabelproperties.libsonnet',
  Regulator: import 'regulator.libsonnet',
  RegulatorDynamicConfig: import 'regulatordynamicconfig.libsonnet',
  RegulatorIns: import 'regulatorins.libsonnet',
  RegulatorParameters: import 'regulatorparameters.libsonnet',
  Resources: import 'resources.libsonnet',
  Rule: import 'rule.libsonnet',
  SMA: import 'sma.libsonnet',
  SMAIns: import 'smains.libsonnet',
  SMAOuts: import 'smaouts.libsonnet',
  SMAParameters: import 'smaparameters.libsonnet',
  ScaleInController: import 'scaleincontroller.libsonnet',
  ScaleInControllerController: import 'scaleincontrollercontroller.libsonnet',
  ScaleOutController: import 'scaleoutcontroller.libsonnet',
  ScaleOutControllerController: import 'scaleoutcontrollercontroller.libsonnet',
  Scheduler: import 'scheduler.libsonnet',
  SchedulerWorkload: import 'schedulerworkload.libsonnet',
  SchedulerWorkloadParameters: import 'schedulerworkloadparameters.libsonnet',
  Selector: import 'selector.libsonnet',
  SignalGenerator: import 'signalgenerator.libsonnet',
  SignalGeneratorIns: import 'signalgeneratorins.libsonnet',
  SignalGeneratorOuts: import 'signalgeneratorouts.libsonnet',
  SignalGeneratorParameters: import 'signalgeneratorparameters.libsonnet',
  SignalGeneratorParametersStep: import 'signalgeneratorparametersstep.libsonnet',
  Switcher: import 'switcher.libsonnet',
  SwitcherIns: import 'switcherins.libsonnet',
  SwitcherOuts: import 'switcherouts.libsonnet',
  TelemetryCollector: import 'telemetrycollector.libsonnet',
  UnaryOperator: import 'unaryoperator.libsonnet',
  UnaryOperatorIns: import 'unaryoperatorins.libsonnet',
  UnaryOperatorOuts: import 'unaryoperatorouts.libsonnet',
  Variable: import 'variable.libsonnet',
  VariableDynamicConfig: import 'variabledynamicconfig.libsonnet',
  VariableOuts: import 'variableouts.libsonnet',
}
