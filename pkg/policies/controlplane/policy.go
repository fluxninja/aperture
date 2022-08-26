package controlplane

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	goObjectHash "github.com/benlaurie/objecthash/go/objecthash"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v2"

	configv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/fluxmeter"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// policyModule returns Fx options of Policy for the Main App.
func policyModule() fx.Option {
	// Circuit module options
	componentFactoryOptions := componentFactoryModule()

	return fx.Options(
		circuitFactoryModule(),
		componentFactoryOptions,
	)
}

// Policy invokes the Circuit runtime at tick frequency.
type Policy struct {
	iface.PolicyBase
	// Metric substitution map
	metricSubMap map[string]*metricSub
	// Circuit
	circuit *runtime.Circuit
	// job group
	circuitJobGroup *jobs.JobGroup
	// Job name
	jobName string
	// Evaluation interval determines how often the Circuit gets executed
	evaluationInterval time.Duration
}

// Make sure Policy complies with PolicyAPI interface.
var _ iface.Policy = (*Policy)(nil)

type metricSub struct {
	metricName    string
	labelMatchers []*labels.Matcher
}

// newPolicyOptions creates a new Policy object and returns its Fx options for the per Policy App.
func newPolicyOptions(
	wrapperMessage *configv1.PolicyWrapper,
) (fx.Option, error) {
	// List of options for the policy.
	policyOptions := []fx.Option{}
	policy, compWithPortsList, partialPolicyOption, err := compilePolicyWrapper(wrapperMessage)
	if err != nil {
		return nil, err
	}
	policyOptions = append(policyOptions, partialPolicyOption)
	policyOptions = append(policyOptions, fx.Supply(
		fx.Annotate(policy, fx.As(new(iface.PolicyRead))),
	))

	// Create circuit
	circuit, circuitOption := runtime.NewCircuitAndOptions(compWithPortsList, policy)
	policyOptions = append(policyOptions, circuitOption)

	policyOptions = append(policyOptions, componentFactoryModuleForPolicyApp(circuit))

	policyOptions = append(policyOptions, fx.Supply(fx.Annotate(circuit, fx.As(new(runtime.CircuitAPI)))))
	policy.circuit = circuit

	return fx.Options(policyOptions...), nil
}

// CompilePolicy takes policyMessage and returns a compiled policy. This is a helper method for standalone consumption of policy compiler.
func CompilePolicy(policyMessage *policylangv1.Policy) ([]runtime.CompiledComponentAndPorts, error) {
	wrapperMessage, err := HashAndPolicyWrap(policyMessage, "DoesNotMatter")
	if err != nil {
		return nil, err
	}
	_, compWithPortsList, _, err := compilePolicyWrapper(wrapperMessage)
	if err != nil {
		return nil, err
	}
	return compWithPortsList, nil
}

// compilePolicyWrapper takes policyProto and returns a compiled policy.
func compilePolicyWrapper(wrapperMessage *configv1.PolicyWrapper) (*Policy, []runtime.CompiledComponentAndPorts, fx.Option, error) {
	if wrapperMessage == nil {
		return nil, nil, nil, fmt.Errorf("nil policy wrapper message")
	}

	policy := &Policy{
		PolicyBase: wrapperMessage,
	}
	policy.metricSubMap = make(map[string]*metricSub)

	// Get Policy Proto
	policyProto := wrapperMessage.GetPolicy()
	if policyProto == nil {
		return nil, nil, nil, fmt.Errorf("nil policy proto")
	}
	// Read evaluation interval
	policy.evaluationInterval = policyProto.EvaluationInterval.AsDuration()

	var fluxMeterOptions []fx.Option
	// Initialize flux meters
	for name, fluxMeterProto := range policyProto.FluxMeters {
		fluxMeterOption, err := fluxmeter.NewFluxMeterOptions(name, fluxMeterProto, policy, policy)
		if err != nil {
			return nil, nil, nil, err
		}
		fluxMeterOptions = append(fluxMeterOptions, fluxMeterOption)
	}

	compWithPortsList, partialCircuitOption, err := compileCircuit(policyProto.Circuit, policy)
	if err != nil {
		return nil, nil, nil, err
	}

	return policy, compWithPortsList, fx.Options(
		fx.Options(fluxMeterOptions...),
		partialCircuitOption,
		fx.Invoke(policy.setup),
	), nil
}

// GetEvaluationInterval returns the ID of the policy.
func (policy *Policy) GetEvaluationInterval() time.Duration {
	return policy.evaluationInterval
}

func (policy *Policy) setup(
	lifecycle fx.Lifecycle,
	circuitJobGroup *jobs.JobGroup,
) error {
	// Job name
	policy.jobName = fmt.Sprintf("Policy-%s", policy.GetPolicyName())
	// Job group
	policy.circuitJobGroup = circuitJobGroup

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// Create a job that runs every tick i.e. evaluation_interval. Set timeout duration to half of evaluation_interval
			job := jobs.BasicJob{
				JobFunc: policy.executeTick,
			}
			job.JobName = policy.jobName
			initialDelay := config.Duration{Duration: durationpb.New(time.Duration(0))}
			executionPeriod := config.Duration{Duration: durationpb.New(policy.evaluationInterval)}
			executionTimeout := config.Duration{Duration: durationpb.New(time.Millisecond * 100)}
			jobConfig := jobs.JobConfig{
				InitiallyHealthy: true,
				InitialDelay:     initialDelay,
				ExecutionPeriod:  executionPeriod,
				ExecutionTimeout: executionTimeout,
			}
			// Register job with registry
			err := policy.circuitJobGroup.RegisterJob(&job, jobConfig)
			if err != nil {
				log.Error().Err(err).Str("job", policy.jobName).Msg("Error registering job")
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			// Deregister job from registry
			_ = policy.circuitJobGroup.DeregisterJob(policy.jobName)
			return nil
		},
	})

	return nil
}

func (policy *Policy) executeTick(jobCtxt context.Context) (proto.Message, error) {
	// Get JobInfo
	jobInfo := policy.circuitJobGroup.JobInfo(policy.jobName)
	if jobInfo == nil {
		return nil, fmt.Errorf("job info not found for job %s", policy.jobName)
	}
	tickInfo := runtime.TickInfo{
		Timestamp:     jobInfo.LastRunTime,
		NextTimestamp: jobInfo.NextRunTime,
		Tick:          jobInfo.RunCount,
		Interval:      policy.evaluationInterval,
	}
	// Execute Circuit
	err := policy.circuit.Execute(tickInfo)
	// TODO: return tick info (publish to health framework) instead of returning nil proto.Message
	return nil, err
}

// RegisterHistogramSub registers a metric substitution pattern for a histogram.
// Not thread safe!
func (policy *Policy) RegisterHistogramSub(metricNameOrig, metricNameSub string, labelMatchers []*labels.Matcher) {
	// Assume histogram metric
	policy.RegisterMetricSub(metricNameOrig+"_sum", metricNameSub+"_sum", labelMatchers)
	policy.RegisterMetricSub(metricNameOrig+"_count", metricNameSub+"_count", labelMatchers)
	policy.RegisterMetricSub(metricNameOrig+"_bucket", metricNameSub+"_bucket", labelMatchers)
}

// RegisterMetricSub registers a metric substitution pattern for a metric.
// Not thread safe!
func (policy *Policy) RegisterMetricSub(metricsNameOrig, metricNameSub string, labelMatchers []*labels.Matcher) {
	metricSub := &metricSub{metricName: metricNameSub, labelMatchers: labelMatchers}
	policy.metricSubMap[metricsNameOrig] = metricSub
}

// ResolveMetricNames resolves metric names based on the substitution patterns.
func (policy *Policy) ResolveMetricNames(query string) (string, error) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return query, err
	}
	msv := &metricSubVisitor{policy: policy}
	var path []parser.Node
	err = parser.Walk(msv, expr, path)
	if err != nil {
		return query, err
	}
	newQuery := expr.String()
	return newQuery, nil
}

func (policy *Policy) resolveMetric(metricName string, labelMatchers []*labels.Matcher) (string, []*labels.Matcher) {
	if metricName != "" {
		metricSub, ok := policy.metricSubMap[metricName]
		if ok {
			// Remove metric label if it exists in labelMatchers to prevent double setting the metric name
			for i, labelMatcher := range labelMatchers {
				if labelMatcher.Name == "__name__" {
					labelMatchers = append(labelMatchers[:i], labelMatchers[i+1:]...)
					break
				}
			}
			return metricSub.metricName, append(labelMatchers, metricSub.labelMatchers...)
		}
	}

	return metricName, labelMatchers
}

type metricSubVisitor struct {
	policy *Policy
}

// Visit implements the visitor interface.
func (msv *metricSubVisitor) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	if node == nil {
		return msv, nil
	}

	switch n := node.(type) {
	case *parser.VectorSelector:
		n.Name, n.LabelMatchers = msv.policy.resolveMetric(n.Name, n.LabelMatchers)
	case *parser.MatrixSelector:
		switch v := n.VectorSelector.(type) {
		case *parser.VectorSelector:
			v.Name, v.LabelMatchers = msv.policy.resolveMetric(v.Name, v.LabelMatchers)
		}
	case *parser.AggregateExpr:
	case *parser.BinaryExpr:
	case *parser.Call:
	case *parser.Expressions:
	case *parser.NumberLiteral:
	case *parser.StringLiteral:
	case *parser.SubqueryExpr:
	case *parser.ParenExpr:
	case *parser.UnaryExpr:
	default:
		log.Warn().Msgf("Unknown type %T", n)
	}
	return msv, nil
}

// HashAndPolicyWrap wraps a proto message with a config properties wrapper and hashes it.
func HashAndPolicyWrap(policyMessage *policylangv1.Policy, policyName string) (*configv1.PolicyWrapper, error) {
	dat, marshalErr := yaml.Marshal(policyMessage)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msgf("Failed to marshal proto message %+v", policyMessage)
		return nil, marshalErr
	}
	hashBytes, hashErr := goObjectHash.ObjectHash(dat)
	if hashErr != nil {
		log.Warn().Err(hashErr).Msgf("Failed to hash json serialized proto message %s", string(dat))
		return nil, hashErr
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes[:])

	return &configv1.PolicyWrapper{
		Policy:     policyMessage,
		PolicyName: policyName,
		PolicyHash: hash,
	}, nil
}
