local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local circuit = spec.v1.Circuit;
local component = spec.v1.Component;
local port = spec.v1.Port;
local combinator = spec.v1.ArithmeticCombinator;
local gradient = spec.v1.GradientController;
local gradientParameters = spec.v1.GradientParameters;
local concurrencyLimiter = spec.v1.ConcurrencyLimiter;
local scheduler = spec.v1.Scheduler;
local schedulerParameters = spec.v1.SchedulerParameters;
local decider = spec.v1.Decider;
local switcher = spec.v1.Switcher;
local loadActuator = spec.v1.LoadActuator;
local loadActuatorDynamicConfig = spec.v1.LoadActuatorDynamicConfig;
local alerterConfig = spec.v1.AlerterConfig;
local max = spec.v1.Max;
local min = spec.v1.Min;
local sqrt = spec.v1.Sqrt;
local firstValid = spec.v1.FirstValid;
local extrapolator = spec.v1.Extrapolator;
local nestedSignalIngress = spec.v1.NestedSignalIngress;
local nestedSignalEgress = spec.v1.NestedSignalEgress;

function(params) {
  _config:: config.circuit + params,

  local c = $._config.constants,

  local circuit =
    [
      component.withArithmeticCombinator(combinator.div(port.withSignalName('DESIRED_CONCURRENCY'),
                                                        port.withSignalName('INCOMING_CONCURRENCY'),
                                                        output=port.withSignalName('DESIRED_CONCURRENCY_RATIO'))),
      component.withArithmeticCombinator(combinator.mul(port.withConstantSignal(c.concurrencyLimitMultiplier),
                                                        port.withSignalName('ACCEPTED_CONCURRENCY'),
                                                        output=port.withSignalName('NORMAL_CONCURRENCY_LIMIT'))),
      component.withArithmeticCombinator(combinator.add(port.withConstantSignal(c.concurrencyLinearIncrement),
                                                        port.withSignalName('SQRT_CONCURRENCY_INCREMENT'),
                                                        output=port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'))),
      component.withArithmeticCombinator(combinator.add(port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'),
                                                        port.withSignalName('CONCURRENCY_INCREMENT'),
                                                        output=port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL'))),
      component.withMin(
        min.new()
        + min.withInPorts({ inputs: [port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL'), port.withSignalName('ACCEPTED_CONCURRENCY')] })
        + min.withOutPorts({ output: port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL_CAPPED') }),
      ),
      component.withFirstValid(
        firstValid.new()
        + firstValid.withInPorts({ inputs: [port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL_CAPPED'), port.withConstantSignal(0)] })
        + firstValid.withOutPorts({ output: port.withSignalName('CONCURRENCY_INCREMENT_NORMAL') }),
      ),
      component.withSqrt(
        sqrt.new()
        + sqrt.withInPorts({ input: port.withSignalName('ACCEPTED_CONCURRENCY') })
        + sqrt.withOutPorts({ output: port.withSignalName('SQRT_CONCURRENCY_INCREMENT') })
        + sqrt.withScale($._config.constants.concurrencySQRTIncrementMultiplier),
      ),
      component.withDecider(
        decider.new()
        + decider.withOperator('gt')
        + decider.withInPortsMixin(
          decider.inPorts.withLhs(port.withSignalName('SIGNAL'))
          + decider.inPorts.withRhs(port.withSignalName('SETPOINT'))
        )
        + decider.withOutPortsMixin(decider.outPorts.withOutput(port.withSignalName('IS_OVERLOAD')))
      ),
      component.withSwitcher(
        switcher.new()
        + switcher.withInPortsMixin(
          switcher.inPorts.withOnTrue(port.withConstantSignal(0))
          + switcher.inPorts.withOnFalse(port.withSignalName('CONCURRENCY_INCREMENT_NORMAL'))
          + switcher.inPorts.withSwitch(port.withSignalName('IS_OVERLOAD'))
        )
        + switcher.withOutPortsMixin(switcher.outPorts.withOutput(port.withSignalName('CONCURRENCY_INCREMENT')))
      ),
      component.withGradientController(
        local g = $._config.gradient;
        gradient.new()
        + gradient.withSlope(g.slope)
        + gradient.withMinGradient(g.minGradient)
        + gradient.withMaxGradient(g.maxGradient)
        + gradient.withInPorts({
          signal: port.withSignalName('SIGNAL'),
          setpoint: port.withSignalName('SETPOINT'),
          max: port.withSignalName('NORMAL_CONCURRENCY_LIMIT'),
          control_variable: port.withSignalName('ACCEPTED_CONCURRENCY'),
          optimize: port.withSignalName('CONCURRENCY_INCREMENT'),
        })
        + gradient.withOutPortsMixin({
          output: port.withSignalName('DESIRED_CONCURRENCY'),
        })
      ),
      component.withExtrapolator(
        extrapolator.new()
        + extrapolator.withMaxExtrapolationInterval('5s')
        + extrapolator.withInPortsMixin({
          input: port.withSignalName('DESIRED_CONCURRENCY_RATIO'),
        })
        + extrapolator.withOutPortsMixin({
          output: port.withSignalName('LOAD_MULTIPLIER'),
        })
      ),
      component.withConcurrencyLimiter(
        local c = $._config.concurrencyLimiter;
        concurrencyLimiter.new()
        + concurrencyLimiter.withFlowSelector($._config.concurrencyLimiterFlowSelector)
        + concurrencyLimiter.withScheduler(
          scheduler.new()
          + scheduler.withAutoTokens(c.autoTokens)
          + scheduler.withTimeoutFactor(c.timeoutFactor)
          + scheduler.withDefaultWorkloadParameters(c.defaultWorkloadParameters)
          + scheduler.withWorkloads(c.workloads)
          + scheduler.withOutPortsMixin({
            accepted_concurrency: port.withSignalName('ACCEPTED_CONCURRENCY'),
            incoming_concurrency: port.withSignalName('INCOMING_CONCURRENCY'),
          })
        )
        + concurrencyLimiter.withLoadActuator(
          loadActuator.withInPortsMixin({ load_multiplier: port.withSignalName('LOAD_MULTIPLIER') })
          + loadActuator.withDynamicConfigKey(c.dynamicConfigKey)
          + loadActuator.withAlerterConfig(
            alerterConfig.new()
            + alerterConfig.withAlertName(c.alerterName)
            + alerterConfig.withAlertChannels(c.alerterChannels)
            + alerterConfig.withResolveTimeout(c.alerterResolveTimeout)
          )
        )
      ),
      component.withNestedSignalIngress(
        nestedSignalIngress.new()
        + nestedSignalIngress.withPortName('signal')
        + nestedSignalIngress.withOutPortsMixin({
          signal: port.withSignalName('SIGNAL'),
        })
      ),
      component.withNestedSignalIngress(
        nestedSignalIngress.new()
        + nestedSignalIngress.withPortName('setpoint')
        + nestedSignalIngress.withOutPortsMixin({
          signal: port.withSignalName('SETPOINT'),
        })
      ),
      component.withNestedSignalEgress(
        nestedSignalEgress.new()
        + nestedSignalEgress.withPortName('accepted_concurrency')
        + nestedSignalEgress.withInPortsMixin({
          signal: port.withSignalName('ACCEPTED_CONCURRENCY'),
        })
      ),
      component.withNestedSignalEgress(
        nestedSignalEgress.new()
        + nestedSignalEgress.withPortName('incoming_concurrency')
        + nestedSignalEgress.withInPortsMixin({
          signal: port.withSignalName('INCOMING_CONCURRENCY'),
        })
      ),
      component.withNestedSignalEgress(
        nestedSignalEgress.new()
        + nestedSignalEgress.withPortName('desired_concurrency')
        + nestedSignalEgress.withInPortsMixin({
          signal: port.withSignalName('DESIRED_CONCURRENCY'),
        })
      ),
      component.withNestedSignalEgress(
        nestedSignalEgress.new()
        + nestedSignalEgress.withPortName('is_overload')
        + nestedSignalEgress.withInPortsMixin({
          signal: port.withSignalName('IS_OVERLOAD'),
        })
      ),
      component.withNestedSignalEgress(
        nestedSignalEgress.new()
        + nestedSignalEgress.withPortName('load_multiplier')
        + nestedSignalEgress.withInPortsMixin({
          signal: port.withSignalName('LOAD_MULTIPLIER'),
        })
      ),
    ] + $._config.components,

  circuit: circuit,
}
