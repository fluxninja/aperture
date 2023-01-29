{
  FluxMeterExponentialBuckets: import 'fluxmeterexponentialbuckets.libsonnet',
  FluxMeterExponentialBucketsRange: import 'fluxmeterexponentialbucketsrange.libsonnet',
  FluxMeterLinearBuckets: import 'fluxmeterlinearbuckets.libsonnet',
  FluxMeterStaticBuckets: import 'fluxmeterstaticbuckets.libsonnet',
  HorizontalPodScalerScaleActuator: import 'horizontalpodscalerscaleactuator.libsonnet',
  HorizontalPodScalerScaleActuatorDynamicConfig: import 'horizontalpodscalerscaleactuatordynamicconfig.libsonnet',
  HorizontalPodScalerScaleActuatorIns: import 'horizontalpodscalerscaleactuatorins.libsonnet',
  HorizontalPodScalerScaleReporter: import 'horizontalpodscalerscalereporter.libsonnet',
  HorizontalPodScalerScaleReporterOuts: import 'horizontalpodscalerscalereporterouts.libsonnet',
  MatchExpressionList: import 'matchexpressionlist.libsonnet',
  RateLimiterLazySync: import 'ratelimiterlazysync.libsonnet',
  RateLimiterOverride: import 'ratelimiteroverride.libsonnet',
  RuleRego: import 'rulerego.libsonnet',
  SchedulerParametersWorkload: import 'schedulerparametersworkload.libsonnet',
  SchedulerParametersWorkloadParameters: import 'schedulerparametersworkloadparameters.libsonnet',
  AIMDConcurrencyController: import 'aimdconcurrencycontroller.libsonnet',
  AIMDConcurrencyControllerIns: import 'aimdconcurrencycontrollerins.libsonnet',
  AIMDConcurrencyControllerOuts: import 'aimdconcurrencycontrollerouts.libsonnet',
  AddressExtractor: import 'addressextractor.libsonnet',
  Alerter: import 'alerter.libsonnet',
  AlerterConfig: import 'alerterconfig.libsonnet',
  AlerterIns: import 'alerterins.libsonnet',
  And: import 'and.libsonnet',
  AndIns: import 'andins.libsonnet',
  AndOuts: import 'andouts.libsonnet',
  ArithmeticCombinator: import 'arithmeticcombinator.libsonnet',
  ArithmeticCombinatorIns: import 'arithmeticcombinatorins.libsonnet',
  ArithmeticCombinatorOuts: import 'arithmeticcombinatorouts.libsonnet',
  Circuit: import 'circuit.libsonnet',
  Classifier: import 'classifier.libsonnet',
  Component: import 'component.libsonnet',
  ConcurrencyLimiter: import 'concurrencylimiter.libsonnet',
  ConstantSignal: import 'constantsignal.libsonnet',
  ControllerDynamicConfig: import 'controllerdynamicconfig.libsonnet',
  Decider: import 'decider.libsonnet',
  DeciderIns: import 'deciderins.libsonnet',
  DeciderOuts: import 'deciderouts.libsonnet',
  Differentiator: import 'differentiator.libsonnet',
  DifferentiatorIns: import 'differentiatorins.libsonnet',
  DifferentiatorOuts: import 'differentiatorouts.libsonnet',
  EMA: import 'ema.libsonnet',
  EMAIns: import 'emains.libsonnet',
  EMAOuts: import 'emaouts.libsonnet',
  EqualsMatchExpression: import 'equalsmatchexpression.libsonnet',
  Extractor: import 'extractor.libsonnet',
  Extrapolator: import 'extrapolator.libsonnet',
  ExtrapolatorIns: import 'extrapolatorins.libsonnet',
  ExtrapolatorOuts: import 'extrapolatorouts.libsonnet',
  FirstValid: import 'firstvalid.libsonnet',
  FirstValidIns: import 'firstvalidins.libsonnet',
  FirstValidOuts: import 'firstvalidouts.libsonnet',
  FlowMatcher: import 'flowmatcher.libsonnet',
  FlowSelector: import 'flowselector.libsonnet',
  FluxMeter: import 'fluxmeter.libsonnet',
  GradientController: import 'gradientcontroller.libsonnet',
  GradientControllerIns: import 'gradientcontrollerins.libsonnet',
  GradientControllerOuts: import 'gradientcontrollerouts.libsonnet',
  GradientParameters: import 'gradientparameters.libsonnet',
  Holder: import 'holder.libsonnet',
  HolderIns: import 'holderins.libsonnet',
  HolderOuts: import 'holderouts.libsonnet',
  HorizontalPodScaler: import 'horizontalpodscaler.libsonnet',
  InPort: import 'inport.libsonnet',
  Integration: import 'integration.libsonnet',
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
  LoadActuator: import 'loadactuator.libsonnet',
  LoadActuatorDynamicConfig: import 'loadactuatordynamicconfig.libsonnet',
  LoadActuatorIns: import 'loadactuatorins.libsonnet',
  MatchExpression: import 'matchexpression.libsonnet',
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
  Policy: import 'policy.libsonnet',
  PromQL: import 'promql.libsonnet',
  PromQLOuts: import 'promqlouts.libsonnet',
  PulseGenerator: import 'pulsegenerator.libsonnet',
  PulseGeneratorOuts: import 'pulsegeneratorouts.libsonnet',
  RateLimiter: import 'ratelimiter.libsonnet',
  RateLimiterDynamicConfig: import 'ratelimiterdynamicconfig.libsonnet',
  RateLimiterIns: import 'ratelimiterins.libsonnet',
  Resources: import 'resources.libsonnet',
  Rule: import 'rule.libsonnet',
  Scheduler: import 'scheduler.libsonnet',
  SchedulerOuts: import 'schedulerouts.libsonnet',
  SchedulerParameters: import 'schedulerparameters.libsonnet',
  ServiceSelector: import 'serviceselector.libsonnet',
  Sqrt: import 'sqrt.libsonnet',
  SqrtIns: import 'sqrtins.libsonnet',
  SqrtOuts: import 'sqrtouts.libsonnet',
  Switcher: import 'switcher.libsonnet',
  SwitcherIns: import 'switcherins.libsonnet',
  SwitcherOuts: import 'switcherouts.libsonnet',
  Variable: import 'variable.libsonnet',
  VariableDynamicConfig: import 'variabledynamicconfig.libsonnet',
  VariableOuts: import 'variableouts.libsonnet',
}
