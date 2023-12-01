package promql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/prometheus"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// ErrNoQueriesReturned is returned when no queries are returned by the policy (initial state).
var ErrNoQueriesReturned = errors.New("no queries returned until now")

const (
	promTimeout = time.Second * 5
)

// Module returns fx options for PromQL in the main app.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(
			provideFxOptionsFunc,
			fx.ResultTags(iface.FxOptionsFuncTag),
		)),
	)
}

func provideFxOptionsFunc(promAPI prometheusv1.API, enforcer *prometheus.PrometheusEnforcer) notifiers.FxOptionsFunc {
	return func(key notifiers.Key, _ config.Unmarshaller, _ status.Registry) (fx.Option, error) {
		return fx.Supply(fx.Annotate(promAPI, fx.As(new(prometheusv1.API))),
			fx.Annotate(enforcer, fx.As(new(*prometheus.PrometheusEnforcer))),
		), nil
	}
}

type scalarResultCallback func(float64, error)

type promResultCallback func(prometheusmodel.Value, error)

type scalarResultBroker struct {
	job   jobs.Job
	err   error
	cb    scalarResultCallback
	query string
	// lock to protect concurrent access to the result
	lock sync.RWMutex
	res  float64
}

// Make sure scalarResultBroker complies with the jobResultBroker interface.
var _ runtime.BackgroundJob = (*scalarResultBroker)(nil)

// GetJob implements runtime.BackgroundJob.GetJob.
func (srb *scalarResultBroker) GetJob() jobs.Job {
	return srb.job
}

// NotifyCompletion implements runtime.BackgroundJob.NotifyCompletion.
func (srb *scalarResultBroker) NotifyCompletion() {
	srb.lock.RLock()
	defer srb.lock.RUnlock()
	srb.cb(srb.res, srb.err)
}

func (srb *scalarResultBroker) handleResult(_ context.Context, value float64, cbArgs ...interface{}) (proto.Message, error) {
	srb.lock.Lock()
	defer srb.lock.Unlock()
	srb.res = value
	srb.err = nil
	return wrapperspb.Double(value), nil
}

func (srb *scalarResultBroker) handleError(err error, cbArgs ...interface{}) (proto.Message, error) {
	srb.lock.Lock()
	defer srb.lock.Unlock()
	srb.res = math.NaN()
	srb.err = err
	return nil, errors.New("invalid ScalarResult, error in prometheus query")
}

type taggedResultBroker struct {
	job   jobs.Job
	res   prometheusmodel.Value
	err   error
	cb    promResultCallback
	query string
	// lock to protect concurrent access to the result
	lock sync.RWMutex
}

// Make sure promResultBroker complies with the jobResultBroker interface.
var _ runtime.BackgroundJob = (*taggedResultBroker)(nil)

// GetJob implements runtime.Job.GetJob.
func (trb *taggedResultBroker) GetJob() jobs.Job {
	return trb.job
}

// NotifyCompletion implements runtime.BackgroundJob.NotifyCompletion.
func (trb *taggedResultBroker) NotifyCompletion() {
	trb.lock.RLock()
	defer trb.lock.RUnlock()
	trb.cb(trb.res, trb.err)
}

func (trb *taggedResultBroker) handleResult(_ context.Context, value prometheusmodel.Value, cbArgs ...interface{}) (proto.Message, error) {
	trb.lock.Lock()
	defer trb.lock.Unlock()
	trb.res = value
	trb.err = nil
	return nil, nil
}

func (trb *taggedResultBroker) handleError(err error, cbArgs ...interface{}) (proto.Message, error) {
	trb.lock.Lock()
	defer trb.lock.Unlock()
	trb.err = err
	trb.res = nil
	return nil, errors.New("invalid taggedResult, error in prometheus query")
}

// Job Register can determine the type of job to register.
type queryScheduerIfc interface {
	scheduleQuery(endTimestamp time.Time, circuitAPI runtime.CircuitAPI)
}

// PromQL is a component that runs a Prometheus query in the background and returns the result as a signal Reading.
type PromQL struct {
	// Prometheus API
	promAPI prometheusv1.API
	// Policy read API
	policyReadAPI iface.Policy
	// Current error
	err error
	// Query Scheduler used to scheduler Prometheus query jobs. Determines the type of job to register.
	queryScheduler queryScheduerIfc
	// Prometheus Labels Enforcer
	enforcer *prometheus.PrometheusEnforcer
	// Job name for the query
	jobName string
	// The query to run
	queryString string
	// Component index
	componentID string
	// Metrics in the query
	metrics []string
	// Execute the query every ticksPerExecution ticks
	ticksPerExecution int
	// Current value
	value float64
	// Interval of time between evaluations
	evaluationInterval time.Duration
}

// Name implements runtime.Component.
func (*PromQL) Name() string { return "PromQL" }

// Type implements runtime.Component.
func (*PromQL) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.
func (promQL *PromQL) ShortDescription() string {
	// print interval and metrics
	metricsList := strings.Join(promQL.metrics, ", ")
	return fmt.Sprintf("runs every %s, metrics: %v", promQL.evaluationInterval, metricsList)
}

// IsActuator implements runtime.Component.
func (*PromQL) IsActuator() bool { return false }

var _ runtime.Component = (*PromQL)(nil)

// Make sure PromQL implements querySchedulerIfc.
var _ queryScheduerIfc = (*PromQL)(nil)

// NewPromQLAndOptions creates PromQL and its fx options.
func NewPromQLAndOptions(
	promQLProto *policylangv1.PromQL,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*PromQL, fx.Option, error) {
	promQLEvaluationInterval := promQLProto.EvaluationInterval.AsDuration()
	circuitEvaluationInterval := policyReadAPI.GetEvaluationInterval()
	ticksPerExecution := int(math.Ceil(float64(promQLEvaluationInterval) / float64(circuitEvaluationInterval)))

	queryString := promQLProto.GetQueryString()
	// Parse metric names in PromQL and save them in an array
	metrics, warnErr := extractMetrics(queryString)
	if warnErr != nil {
		log.Warn().Err(warnErr).Msg("Could not parse metrics from PromQL query")
	}

	promQL := &PromQL{
		ticksPerExecution:  ticksPerExecution,
		policyReadAPI:      policyReadAPI,
		componentID:        componentID.String(),
		evaluationInterval: circuitEvaluationInterval * time.Duration(ticksPerExecution),
		queryString:        queryString,
		metrics:            metrics,
		// Set err to make sure the initial runs of Execute return Invalid readings.
		err: ErrNoQueriesReturned,
	}

	// Job register is implemented by self
	promQL.queryScheduler = promQL

	// Job name
	promQL.jobName = fmt.Sprintf("Component-%s", promQL.componentID)

	// Invoke setup in the Policy app startup via fx.Options
	options := fx.Options(
		fx.Invoke(
			promQL.setup,
		),
	)
	return promQL, options, nil
}

func (promQL *PromQL) setup(promAPI prometheusv1.API, enforcer *prometheus.PrometheusEnforcer) error {
	promQL.promAPI = promAPI
	promQL.enforcer = enforcer

	return nil
}

// Execute implements runtime.Component.Execute.
func (promQL *PromQL) Execute(inPortReadings runtime.PortToReading,
	circuitAPI runtime.CircuitAPI,
) (outPortReadings runtime.PortToReading, err error) {
	tickInfo := circuitAPI.GetTickInfo()
	endTimestamp := tickInfo.Timestamp().Truncate(tickInfo.Interval())
	promQL.scheduleQuery(endTimestamp, circuitAPI)

	// Create current reading based on err and value
	var currentReading runtime.Reading
	if promQL.err != nil {
		currentReading = runtime.InvalidReading()
	} else {
		currentReading = runtime.NewReading(promQL.value)
	}

	return runtime.PortToReading{
		"output": []runtime.Reading{currentReading},
	}, nil
}

// DynamicConfigUpdate is a no-op for PromQL.
func (promQL *PromQL) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}

func (promQL *PromQL) scheduleQuery(endTimestamp time.Time, circuitAPI runtime.CircuitAPI) {
	jobName := promQL.jobName
	query := promQL.queryString
	promAPI := promQL.promAPI
	enforcer := promQL.enforcer
	cb := promQL.onScalarResult
	// Result handler for this job
	scalarResBroker := &scalarResultBroker{
		cb:    cb,
		query: query,
	}
	job := jobs.NewBasicJob(jobName,
		prometheus.NewScalarQueryJob(
			query,
			endTimestamp,
			promAPI,
			enforcer,
			promTimeout,
			scalarResBroker.handleResult,
			scalarResBroker.handleError,
		))
	scalarResBroker.job = job
	circuitAPI.ScheduleConditionalBackgroundJob(scalarResBroker, promQL.ticksPerExecution)
}

func (promQL *PromQL) onScalarResult(value float64, err error) {
	promQL.value = value
	promQL.err = err
}

// ScalarQuery is a construct that can be used by other components to get tick aligned scalar results of a PromQL query.
type ScalarQuery struct {
	promQL *PromQL
}

// NewScalarQueryAndOptions creates a new ScalarQuery and its fx options.
func NewScalarQueryAndOptions(
	queryString string,
	evaluationInterval time.Duration,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
	jobPostFix string,
) (*ScalarQuery, fx.Option, error) {
	// Create promQLProto
	promQLProto := &policylangv1.PromQL{
		QueryString:        queryString,
		EvaluationInterval: durationpb.New(evaluationInterval),
	}
	// Create promQL
	promQL, options, err := NewPromQLAndOptions(promQLProto, componentID, policyReadAPI)
	if err != nil {
		return nil, fx.Options(), err
	}
	promQL.jobName = fmt.Sprintf("Component-%s.%s", promQL.componentID, jobPostFix)
	scalarQuery := &ScalarQuery{
		promQL: promQL,
	}
	return scalarQuery, options, nil
}

// ExecuteScalarQuery runs a ScalarQueryJob and returns the current results: value and err. This function is supposed to be run under Circuit Execution Lock (Execution of Circuit Components is protected by this lock).
func (scalarQuery *ScalarQuery) ExecuteScalarQuery(circuitAPI runtime.CircuitAPI) (ScalarResult, error) {
	inPortReadings := runtime.PortToReading{}
	_, _ = scalarQuery.promQL.Execute(inPortReadings, circuitAPI)
	// FYI: promQL ensures that initial runs return err when no queries have returned yet.
	return ScalarResult{Value: scalarQuery.promQL.value}, scalarQuery.promQL.err
}

// ScalarResult is the result of a ScalarQuery.
type ScalarResult struct {
	Value float64
}

// TaggedQuery is a construct that can be used by other components to get tick aligned prometheus value results of a PromQL query.
type TaggedQuery struct {
	scalarQuery *ScalarQuery
	res         prometheusmodel.Value
	err         error
}

// Make sure TaggedQuery implements queryScheduerIfc.
var _ queryScheduerIfc = (*TaggedQuery)(nil)

// NewTaggedQueryAndOptions creates a new TaggedQuery and its fx options.
func NewTaggedQueryAndOptions(
	queryString string,
	evaluationInterval time.Duration,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
	jobPostFix string,
) (*TaggedQuery, fx.Option, error) {
	scalarQuery, options, err := NewScalarQueryAndOptions(queryString, evaluationInterval, componentID, policyReadAPI, jobPostFix)
	if err != nil {
		return nil, fx.Options(), err
	}
	taggedQuery := &TaggedQuery{
		scalarQuery: scalarQuery,
		// Set err to make sure the initial runs of ExecuteTaggedQuery return error.
		err: ErrNoQueriesReturned,
	}
	// taggedQuery implements jobRegisterer
	scalarQuery.promQL.queryScheduler = taggedQuery
	return taggedQuery, options, nil
}

func (taggedQuery *TaggedQuery) scheduleQuery(endTimestamp time.Time, circuitAPI runtime.CircuitAPI) {
	jobName := taggedQuery.scalarQuery.promQL.jobName
	query := taggedQuery.scalarQuery.promQL.queryString
	promAPI := taggedQuery.scalarQuery.promQL.promAPI
	enforcer := taggedQuery.scalarQuery.promQL.enforcer
	cb := taggedQuery.onTaggedResult
	// Result handler for this job
	taggedResBroker := &taggedResultBroker{
		cb:    cb,
		query: query,
	}
	job := jobs.NewBasicJob(jobName,
		prometheus.NewPromQueryJob(
			query,
			endTimestamp,
			promAPI,
			enforcer,
			promTimeout,
			taggedResBroker.handleResult,
			taggedResBroker.handleError,
		))
	taggedResBroker.job = job
	circuitAPI.ScheduleConditionalBackgroundJob(taggedResBroker, taggedQuery.scalarQuery.promQL.ticksPerExecution)
}

func (taggedQuery *TaggedQuery) onTaggedResult(res prometheusmodel.Value, err error) {
	taggedQuery.res = res
	taggedQuery.err = err
}

// ExecuteTaggedQuery runs a PromQueryJob and returns the current results: res and err. This function is supposed to be run under Circuit Execution Lock (Execution of Circuit Components is protected by this lock).
func (taggedQuery *TaggedQuery) ExecuteTaggedQuery(circuitAPI runtime.CircuitAPI) (TaggedResult, error) {
	_, _ = taggedQuery.scalarQuery.ExecuteScalarQuery(circuitAPI)
	return TaggedResult{Value: taggedQuery.res}, taggedQuery.err
}

// TaggedResult is the result of a ScalarQuery.
type TaggedResult struct {
	Value prometheusmodel.Value
}

func extractMetrics(query string) ([]string, error) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]struct{})

	// Walk through the PromQL expression and extract metric names
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		if n, ok := node.(*parser.VectorSelector); ok {
			metrics[n.Name] = struct{}{}
		}
		return nil
	})

	// Convert map keys to slice
	var result []string
	for metric := range metrics {
		result = append(result, metric)
	}

	return result, nil
}
