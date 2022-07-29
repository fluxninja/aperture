package controlplane

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	configv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/config/v1"
	policylangv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/FluxNinja/aperture/pkg/config"
	etcdclient "github.com/FluxNinja/aperture/pkg/etcd/client"
	"github.com/FluxNinja/aperture/pkg/jobs"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/policies/apis/policyapi"
	"github.com/FluxNinja/aperture/pkg/policies/controlplane/fluxmeter"
	"github.com/FluxNinja/aperture/pkg/policies/controlplane/runtime"
)

// PolicyModule returns Fx options of Policy for the Main App.
func PolicyModule() fx.Option {
	// Circuit module options
	componentFactoryOptions := ComponentFactoryModule()

	return fx.Options(
		CircuitFactoryModule(),
		componentFactoryOptions,
	)
}

// Policy invokes the Circuit runtime at tick frequency.
type Policy struct {
	policyapi.PolicyBaseAPI
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
var _ policyapi.PolicyAPI = (*Policy)(nil)

type metricSub struct {
	metricName    string
	labelMatchers []*labels.Matcher
}

// NewPolicyOptions creates a new Policy object and returns its Fx options for the per Policy App.
func NewPolicyOptions(
	circuitJobGroup *jobs.JobGroup,
	etcdClient *etcdclient.Client,
	wrapperMessage *configv1.ConfigPropertiesWrapper,
	policyProto *policylangv1.Policy,
) (fx.Option, error) {
	policy := &Policy{
		PolicyBaseAPI:   wrapperMessage,
		circuitJobGroup: circuitJobGroup,
	}

	policy.metricSubMap = make(map[string]*metricSub)

	// Read evaluation interval
	policy.evaluationInterval = policyProto.EvaluationInterval.AsDuration()

	var fluxMeterOptions []fx.Option
	// Initialize flux meters
	for _, fluxMeterProto := range policyProto.FluxMeters {
		fluxMeterOption, err := fluxmeter.NewFluxMeterOptions(fluxMeterProto, policy, policy)
		if err != nil {
			return nil, err
		}
		fluxMeterOptions = append(fluxMeterOptions, fluxMeterOption)
	}

	// Initialize circuit
	circuit, circuitOptions, err := NewCircuitAndOptions(policyProto.Circuit, policy)
	if err != nil {
		return nil, err
	}
	policy.circuit = circuit

	return fx.Options(
		fx.Supply(
			fx.Annotate(policy, fx.As(new(policyapi.PolicyReadAPI))),
			etcdClient,
		),
		circuitOptions,
		fx.Options(fluxMeterOptions...),
		fx.Invoke(policy.setup),
	), nil
}

// GetEvaluationInterval returns the ID of the policy.
func (policy *Policy) GetEvaluationInterval() time.Duration {
	return policy.evaluationInterval
}

func (policy *Policy) setup(lifecycle fx.Lifecycle) error {
	// Job name
	policy.jobName = fmt.Sprintf("Policy-%s", policy.GetPolicyName())

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
