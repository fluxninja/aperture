package promql

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
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

var promQLJobGroupTag = iface.PoliciesRoot + "promql_jobs"

// Module returns fx options for PromQL in the main app.
func Module() fx.Option {
	return fx.Options(
		jobs.JobGroupConstructor{Name: promQLJobGroupTag, Key: iface.PoliciesRoot + ".promql_jobs_scheduler"}.Annotate(),
		fx.Provide(fx.Annotate(
			provideFxOptionsFunc,

			fx.ParamTags(config.NameTag(promQLJobGroupTag)),
			fx.ResultTags(iface.FxOptionsFuncTag),
		)),
	)
}

func provideFxOptionsFunc(promQLJobGroup *jobs.JobGroup, promAPI prometheusv1.API, enforcer *prometheus.PrometheusEnforcer) notifiers.FxOptionsFunc {
	return func(key notifiers.Key, _ config.Unmarshaller, _ status.Registry) (fx.Option, error) {
		return fx.Supply(fx.Annotated{Name: promQLJobGroupTag, Target: promQLJobGroup},
			fx.Annotate(promAPI, fx.As(new(prometheusv1.API))),
			fx.Annotate(enforcer, fx.As(new(*prometheus.PrometheusEnforcer))),
		), nil
	}
}

// ModuleForPolicyApp returns fx options for PromQL in the policy app. Invoked only once per policy.
func ModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	providePromJobsExecutor := func(promQLJobGroup *jobs.JobGroup, lifecycle fx.Lifecycle) (*promJobsExecutor, error) {
		// Create this watcher as a singleton at the policy/circuit level
		pje := &promJobsExecutor{
			circuitAPI:     circuitAPI,
			inflightJobs:   make(jobResultBrokers),
			pendingJobs:    make(jobResultBrokers),
			promQLJobGroup: promQLJobGroup,
		}
		// Register TickEndCallback
		circuitAPI.RegisterTickEndCallback(pje.onTickEnd)

		var jws []jobs.JobWatcher
		jws = append(jws, pje)

		// Create promMultiJob for this circuit
		promMultiJob := jobs.NewMultiJob(promQLJobGroup.GetStatusRegistry().Child(
			fmt.Sprintf("%s-promql", circuitAPI.GetPolicyHash()), circuitAPI.GetPolicyName()), jws, nil)
		pje.promMultiJob = promMultiJob

		executionPeriod := config.MakeDuration(-1)
		executionTimeout := config.MakeDuration(promTimeout * 2)
		jobConfig := jobs.JobConfig{
			InitiallyHealthy: true,
			ExecutionPeriod:  executionPeriod,
			ExecutionTimeout: executionTimeout,
		}

		// Lifecycle hooks to register and deregister this circuit is promMultiJob from promQLJobGroup
		lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				// Register multi job with job group
				err := promQLJobGroup.RegisterJob(promMultiJob, jobConfig)
				return err
			},
			OnStop: func(_ context.Context) error {
				// Deregister multi job from job group
				err := promQLJobGroup.DeregisterJob(promMultiJob.Name())
				return err
			},
		})
		return pje, nil
	}

	return fx.Options(
		fx.Provide(fx.Annotate(
			providePromJobsExecutor,
			fx.ParamTags(config.NameTag(promQLJobGroupTag)),
		)),
	)
}

type jobResultBrokers map[string]jobResultBroker

type scalarResultCallback func(float64, error)

type promResultCallback func(prometheusmodel.Value, error)

type promJobsExecutor struct {
	// CircuitAPI
	circuitAPI runtime.CircuitAPI
	// inflightJobs contains a Job Result Broker for each job in the multi job
	inflightJobs jobResultBrokers
	// pendingJobs contains a Job Result Broker for each job in the multi job
	pendingJobs jobResultBrokers
	// Prom Multi Job
	promMultiJob *jobs.MultiJob
	// Job group
	promQLJobGroup *jobs.JobGroup
	// Query job state
	jobRunning bool
}

// Make sure promJobsWatcher complies with the jobs.JobsWatcher interface.
var _ jobs.JobWatcher = (*promJobsExecutor)(nil)

func (pje *promJobsExecutor) registerScalarJob(
	jobName string,
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	enforcer *prometheus.PrometheusEnforcer,
	timeout time.Duration,
	cb scalarResultCallback,
) {
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
			timeout,
			scalarResBroker.handleResult,
			scalarResBroker.handleError,
		))
	scalarResBroker.job = job
	pje.pendingJobs[jobName] = scalarResBroker
}

func (pje *promJobsExecutor) registerTaggedJob(
	jobName string,
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	enforcer *prometheus.PrometheusEnforcer,
	timeout time.Duration,
	cb promResultCallback,
) {
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
			timeout,
			taggedResBroker.handleResult,
			taggedResBroker.handleError,
		))
	taggedResBroker.job = job
	pje.pendingJobs[jobName] = taggedResBroker
}

// OnJobScheduled is called when the pje.promMultiJob is scheduled.
func (pje *promJobsExecutor) OnJobScheduled() {
}

// OnJobCompleted is called when the pje.promMultiJob is completed.
func (pje *promJobsExecutor) OnJobCompleted(_ *statusv1.Status, _ jobs.JobStats) {
	// Take circuit execution lock
	pje.circuitAPI.LockExecution()
	defer pje.circuitAPI.UnlockExecution()

	// Provide results via callbacks
	for _, jobResBroker := range pje.inflightJobs {
		jobResBroker.deliverResult()
	}
	// Reset inflightJobs
	pje.inflightJobs = make(jobResultBrokers)
	pje.jobRunning = false
}

func (pje *promJobsExecutor) onTickEnd(_ runtime.TickInfo) (err error) {
	logger := pje.circuitAPI.GetStatusRegistry().GetLogger()
	// Already under circuit execution lock
	// Launch job only if previous one is completed
	if pje.jobRunning {
		err = errors.New("previous job is still running")
	} else {
		pje.jobRunning = true
		// Remove all the previous jobs in the multi job
		pje.promMultiJob.DeregisterAll()
		// Add all the pendingJobs to the multijob and trigger it
		for _, jobResBroker := range pje.pendingJobs {
			job := jobResBroker.getJob()
			err = pje.promMultiJob.RegisterJob(job)
			if err != nil {
				logger.Error().Err(err).Str("job", job.Name()).Msg("Error registering job")
				return err
			}
		}
		// Move pendingJobs to inflightJobs
		pje.inflightJobs = pje.pendingJobs
		// Clear pendingJobs for future ticks
		pje.pendingJobs = make(jobResultBrokers)
		// Trigger the multi job
		pje.promQLJobGroup.TriggerJob(pje.promMultiJob.Name(), time.Duration(0))
	}
	return err
}

type jobResultBroker interface {
	getJob() jobs.Job
	getQuery() string
	deliverResult()
}

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
var _ jobResultBroker = (*scalarResultBroker)(nil)

func (srb *scalarResultBroker) getJob() jobs.Job {
	return srb.job
}

func (srb *scalarResultBroker) getQuery() string {
	return srb.query
}

func (srb *scalarResultBroker) deliverResult() {
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
	cb  promResultCallback
	job jobs.Job
	// lock to protect concurrent access to the result
	lock  sync.RWMutex
	res   prometheusmodel.Value
	err   error
	query string
}

// Make sure promResultBroker complies with the jobResultBroker interface.
var _ jobResultBroker = (*taggedResultBroker)(nil)

func (prb *taggedResultBroker) getJob() jobs.Job {
	return prb.job
}

func (prb *taggedResultBroker) getQuery() string {
	return prb.query
}

func (prb *taggedResultBroker) deliverResult() {
	prb.lock.RLock()
	defer prb.lock.RUnlock()
	prb.cb(prb.res, prb.err)
}

func (prb *taggedResultBroker) handleResult(_ context.Context, value prometheusmodel.Value, cbArgs ...interface{}) (proto.Message, error) {
	prb.lock.Lock()
	defer prb.lock.Unlock()
	prb.res = value
	prb.err = nil
	return nil, nil
}

func (prb *taggedResultBroker) handleError(err error, cbArgs ...interface{}) (proto.Message, error) {
	prb.lock.Lock()
	defer prb.lock.Unlock()
	prb.err = err
	prb.res = nil
	return nil, errors.New("invalid taggedResult, error in prometheus query")
}

// Job Register can determine the type of job to register.
type jobRegistererIfc interface {
	registerJob(endTimestamp time.Time)
}

// PromQL is a component that runs a Prometheus query in the background and returns the result as a signal Reading.
type PromQL struct {
	// Last Query Tick
	tickInfo runtime.TickInfo
	// Prometheus API
	promAPI prometheusv1.API
	// Prometheus Labels Enforcer
	enforcer *prometheus.PrometheusEnforcer
	// Policy read API
	policyReadAPI iface.Policy
	// Current error
	err error
	// Job Registerer used to register Prometheus query jobs. Determines the type of job to register.
	jobRegisterer jobRegistererIfc
	// Job executor
	jobExecutor *promJobsExecutor
	// Job name for the query
	jobName string
	// The query to run
	queryString string
	// Component index
	componentID string
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
	return fmt.Sprintf("every %s", promQL.evaluationInterval)
}

// IsActuator implements runtime.Component.
func (*PromQL) IsActuator() bool { return false }

var _ runtime.Component = (*PromQL)(nil)

// Make sure PromQL implements jobRegistererIfc.
var _ jobRegistererIfc = (*PromQL)(nil)

// NewPromQLAndOptions creates PromQL and its fx options.
func NewPromQLAndOptions(
	promQLProto *policylangv1.PromQL,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (*PromQL, fx.Option, error) {
	promQL := &PromQL{
		evaluationInterval: promQLProto.EvaluationInterval.AsDuration(),
		policyReadAPI:      policyReadAPI,
		componentID:        componentID.String(),
		// Set err to make sure the initial runs of Execute return Invalid readings.
		err: ErrNoQueriesReturned,
	}

	// Job register is implemented by self
	promQL.jobRegisterer = promQL

	// Job name
	promQL.jobName = fmt.Sprintf("Component-%s", promQL.componentID)

	// Resolve metric names in PromQL to get the query string
	promQL.queryString = promQLProto.GetQueryString()

	// Invoke setup in the Policy app startup via fx.Options
	options := fx.Options(
		fx.Invoke(
			promQL.setup,
		),
	)
	return promQL, options, nil
}

func (promQL *PromQL) setup(pje *promJobsExecutor, promAPI prometheusv1.API, enforcer *prometheus.PrometheusEnforcer) error {
	promQL.promAPI = promAPI
	promQL.enforcer = enforcer
	promQL.jobExecutor = pje

	return nil
}

// Execute implements runtime.Component.Execute.
func (promQL *PromQL) Execute(inPortReadings runtime.PortToReading,
	tickInfo runtime.TickInfo,
) (outPortReadings runtime.PortToReading, err error) {
	// Re-run query if evaluationInterval elapsed since last query
	if tickInfo.Timestamp().Sub(promQL.lastQueryTimestamp()) >= promQL.evaluationInterval {
		// Run query
		promQL.tickInfo = tickInfo
		// Launch job only if previous one is completed
		// Quantize endTimestamp of query based on tick interval
		endTimestamp := tickInfo.Timestamp().Truncate(tickInfo.Interval())
		// Register jobFunc with jobExecutor
		promQL.jobRegisterer.registerJob(endTimestamp)
	}

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

func (promQL *PromQL) registerJob(endTimestamp time.Time) {
	promQL.jobExecutor.registerScalarJob(
		promQL.jobName,
		promQL.queryString,
		endTimestamp,
		promQL.promAPI,
		promQL.enforcer,
		promTimeout,
		promQL.onScalarResult,
	)
}

func (promQL *PromQL) onScalarResult(value float64, err error) {
	promQL.value = value
	promQL.err = err
}

func (promQL *PromQL) lastQueryTimestamp() time.Time {
	if promQL.tickInfo == nil {
		return time.Time{}
	}
	return promQL.tickInfo.Timestamp()
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
func (scalarQuery *ScalarQuery) ExecuteScalarQuery(tickInfo runtime.TickInfo) (ScalarResult, error) {
	inPortReadings := runtime.PortToReading{}
	_, _ = scalarQuery.promQL.Execute(inPortReadings, tickInfo)
	// FYI: promQL ensures that initial runs return err when no queries have returned yet.
	return ScalarResult{Value: scalarQuery.promQL.value, TickInfo: scalarQuery.promQL.tickInfo}, scalarQuery.promQL.err
}

// ScalarResult is the result of a ScalarQuery.
type ScalarResult struct {
	TickInfo runtime.TickInfo
	Value    float64
}

// TaggedQuery is a construct that can be used by other components to get tick aligned prometheus value results of a PromQL query.
type TaggedQuery struct {
	scalarQuery *ScalarQuery
	res         prometheusmodel.Value
	err         error
}

// Make sure TaggedQuery implements jobRegistererIfc.
var _ jobRegistererIfc = (*TaggedQuery)(nil)

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
	scalarQuery.promQL.jobRegisterer = taggedQuery
	return taggedQuery, options, nil
}

func (taggedQuery *TaggedQuery) registerJob(endTimestamp time.Time) {
	taggedQuery.scalarQuery.promQL.jobExecutor.registerTaggedJob(
		taggedQuery.scalarQuery.promQL.jobName,
		taggedQuery.scalarQuery.promQL.queryString,
		endTimestamp,
		taggedQuery.scalarQuery.promQL.promAPI,
		taggedQuery.scalarQuery.promQL.enforcer,
		promTimeout,
		taggedQuery.onTaggedResult,
	)
}

func (taggedQuery *TaggedQuery) onTaggedResult(res prometheusmodel.Value, err error) {
	taggedQuery.res = res
	taggedQuery.err = err
}

// ExecuteTaggedQuery runs a PromQueryJob and returns the current results: res and err. This function is supposed to be run under Circuit Execution Lock (Execution of Circuit Components is protected by this lock).
func (taggedQuery *TaggedQuery) ExecuteTaggedQuery(tickInfo runtime.TickInfo) (TaggedResult, error) {
	_, _ = taggedQuery.scalarQuery.ExecuteScalarQuery(tickInfo)
	return TaggedResult{Value: taggedQuery.res, TickInfo: taggedQuery.scalarQuery.promQL.tickInfo}, taggedQuery.err
}

// TaggedResult is the result of a ScalarQuery.
type TaggedResult struct {
	TickInfo runtime.TickInfo
	Value    prometheusmodel.Value
}
