package components

import (
	"context"
	"errors"
	"fmt"
	"time"

	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/common"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/prometheus"
	"github.com/fluxninja/aperture/pkg/status"
)

var errNoQueriesReturned = errors.New("no queries returned until now")

const (
	promTimeout = time.Second * 5
)

var promQLJobGroupTag = iface.PoliciesRoot + "promql_jobs"

// PromQLModule returns fx options for PromQL in the main app.
func PromQLModule() fx.Option {
	return fx.Options(
		jobs.JobGroupConstructor{Name: promQLJobGroupTag, Key: iface.PoliciesRoot + ".promql_jobs_scheduler"}.Annotate(),
		fx.Provide(fx.Annotate(
			provideFxOptionsFunc,
			fx.ParamTags(config.NameTag(promQLJobGroupTag)),
			fx.ResultTags(common.FxOptionsFuncTag),
		)),
	)
}

func provideFxOptionsFunc(promQLJobGroup *jobs.JobGroup, promAPI prometheusv1.API) notifiers.FxOptionsFunc {
	return func(key notifiers.Key, unmarshaller config.Unmarshaller, _ status.Registry) (fx.Option, error) {
		return fx.Supply(fx.Annotated{Name: promQLJobGroupTag, Target: promQLJobGroup},
			fx.Annotate(promAPI, fx.As(new(prometheusv1.API))),
		), nil
	}
}

// PromQLModuleForPolicyApp returns fx options for PromQL in the policy app. Invoked only once per policy.
func PromQLModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	providePromJobsExecutor := func(promQLJobGroup *jobs.JobGroup, lifecycle fx.Lifecycle) (*promJobsExecutor, error) {
		// Create this watcher as a singleton at the policy/circuit level
		pje := &promJobsExecutor{
			circuitAPI:     circuitAPI,
			jobResBrokers:  make(jobResultBrokers),
			promQLJobGroup: promQLJobGroup,
		}
		// Register TickEndCallback
		circuitAPI.RegisterTickEndCallback(pje.onTickEnd)

		var jws []jobs.JobWatcher
		jws = append(jws, pje)

		// Create promMultiJob for this circuit
		promMultiJob := jobs.NewMultiJob(circuitAPI.GetPolicyName(), false, jws, nil)
		pje.promMultiJob = promMultiJob

		initialDelay := config.MakeDuration(-1)
		executionPeriod := config.MakeDuration(-1)
		executionTimeout := config.MakeDuration(promTimeout * 2)
		jobConfig := jobs.JobConfig{
			InitiallyHealthy: true,
			InitialDelay:     initialDelay,
			ExecutionPeriod:  executionPeriod,
			ExecutionTimeout: executionTimeout,
		}

		// Lifecycle hooks to register and deregister this circuit's promMultiJob from promQLJobGroup
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
	// jobResBrokers contains a Job Result Broker for each job in the multi job
	jobResBrokers jobResultBrokers
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
	timeout time.Duration,
	cb scalarResultCallback,
) {
	// Result handler for this job
	scalarResBroker := &scalarResultBroker{
		cb:    cb,
		query: query,
	}
	job := &jobs.BasicJob{
		JobFunc: prometheus.NewScalarQueryJob(
			query,
			endTimestamp,
			promAPI,
			timeout,
			scalarResBroker.handleResult,
			scalarResBroker.handleError,
		),
	}
	job.JobName = jobName
	scalarResBroker.job = job
	pje.jobResBrokers[jobName] = scalarResBroker
}

func (pje *promJobsExecutor) registerTaggedJob(
	jobName string,
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	timeout time.Duration,
	cb promResultCallback,
) {
	// Result handler for this job
	taggedResBroker := &taggedResultBroker{
		cb:    cb,
		query: query,
	}
	job := &jobs.BasicJob{
		JobFunc: prometheus.NewPromQueryJob(
			query,
			endTimestamp,
			promAPI,
			timeout,
			taggedResBroker.handleResult,
			taggedResBroker.handleError,
		),
	}
	job.JobName = jobName
	taggedResBroker.job = job
	pje.jobResBrokers[jobName] = taggedResBroker
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
	for _, jobResBroker := range pje.jobResBrokers {
		jobResBroker.deliverResult()
	}
	pje.jobRunning = false
}

func (pje *promJobsExecutor) onTickEnd(_ runtime.TickInfo) (err error) {
	// Already under circuit execution lock
	// Launch job only if previous one is completed
	if pje.jobRunning {
		err = errors.New("previous job is still running")
	} else {
		pje.jobRunning = true
		// Remove all the previous jobs in the multi job
		pje.promMultiJob.DeregisterAll()
		// Add all the new jobs to the multijob and trigger it
		for _, jobResBroker := range pje.jobResBrokers {
			job := jobResBroker.getJob()
			err = pje.promMultiJob.RegisterJob(job)
			if err != nil {
				log.Error().Err(err).Str("job", job.Name()).Msg("Error registering job")
				return err
			}
		}
		// Clear jobsToRun for the next tick
		pje.jobResBrokers = make(jobResultBrokers)
		// Trigger the multi job
		pje.promQLJobGroup.TriggerJob(pje.promMultiJob.Name())
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
	res   float64
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
	srb.cb(srb.res, srb.err)
}

func (srb *scalarResultBroker) handleResult(_ context.Context, value float64, cbArgs ...interface{}) (proto.Message, error) {
	srb.res = value
	srb.cb(value, nil)
	return wrapperspb.Double(value), nil
}

func (srb *scalarResultBroker) handleError(err error, cbArgs ...interface{}) (proto.Message, error) {
	srb.err = err
	srb.cb(0, err)
	return nil, err
}

type taggedResultBroker struct {
	cb    promResultCallback
	job   jobs.Job
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
	prb.cb(prb.res, prb.err)
}

func (prb *taggedResultBroker) handleResult(_ context.Context, value prometheusmodel.Value, cbArgs ...interface{}) (proto.Message, error) {
	prb.res = value
	prb.cb(value, nil)
	return nil, nil
}

func (prb *taggedResultBroker) handleError(err error, cbArgs ...interface{}) (proto.Message, error) {
	prb.err = err
	prb.cb(nil, err)
	return nil, err
}

// Job Register can determine the type of job to register.
type jobRegistererIfc interface {
	registerJob(endTimestamp time.Time)
}

// PromQL is a component that runs a Prometheus query in the background and returns the result as a signal Reading.
type PromQL struct {
	// Last Query Timestamp
	lastQueryTimestamp time.Time
	// Prometheus API
	promAPI prometheusv1.API
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
	componentIndex int
	// Current value
	value float64
	// Interval of time between evaluations
	evaluationInterval time.Duration
}

// Make sure PromQL implements jobRegistererIfc.
var _ jobRegistererIfc = (*PromQL)(nil)

// NewPromQLAndOptions creates PromQL and its fx options.
func NewPromQLAndOptions(
	promQLProto *policylangv1.PromQL,
	componentIndex int,
	policyReadAPI iface.Policy,
) (*PromQL, fx.Option, error) {
	promQL := &PromQL{
		evaluationInterval: promQLProto.EvaluationInterval.AsDuration(),
		policyReadAPI:      policyReadAPI,
		componentIndex:     componentIndex,
		lastQueryTimestamp: time.Time{},
		// Set err to make sure the initial runs of Execute return Invalid readings.
		err: errNoQueriesReturned,
	}

	// Job register is implemented by self
	promQL.jobRegisterer = promQL

	// Job name
	promQL.jobName = fmt.Sprintf("Component-%d", promQL.componentIndex)

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

func (promQL *PromQL) setup(pje *promJobsExecutor, promAPI prometheusv1.API) error {
	promQL.promAPI = promAPI
	promQL.jobExecutor = pje

	return nil
}

// Execute implements runtime.Component.Execute.
func (promQL *PromQL) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (outPortReadings runtime.PortToValue, err error) {
	// Re-run query if evaluationInterval elapsed since last query
	if tickInfo.Timestamp().Sub(promQL.lastQueryTimestamp) >= promQL.evaluationInterval {
		// Run query
		promQL.lastQueryTimestamp = tickInfo.Timestamp()
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

	return runtime.PortToValue{
		"output": []runtime.Reading{currentReading},
	}, nil
}

func (promQL *PromQL) registerJob(endTimestamp time.Time) {
	promQL.jobExecutor.registerScalarJob(
		promQL.jobName,
		promQL.queryString,
		endTimestamp,
		promQL.promAPI,
		promTimeout,
		promQL.onScalarResult,
	)
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
	componentIndex int,
	policyReadAPI iface.Policy,
	jobPostFix string,
) (*ScalarQuery, fx.Option, error) {
	// Create promQLProto
	promQLProto := &policylangv1.PromQL{
		QueryString:        queryString,
		EvaluationInterval: durationpb.New(evaluationInterval),
	}
	// Create promQL
	promQL, options, err := NewPromQLAndOptions(promQLProto, componentIndex, policyReadAPI)
	if err != nil {
		return nil, fx.Options(), err
	}
	promQL.jobName = fmt.Sprintf("Component-%d.%s", promQL.componentIndex, jobPostFix)
	scalarQuery := &ScalarQuery{
		promQL: promQL,
	}
	return scalarQuery, options, nil
}

// ExecuteScalarQuery runs a ScalarQueryJob and returns the current results: value and err. This function is supposed to be run under Circuit Execution Lock (Execution of Circuit Components is protected by this lock).
func (scalarQuery *ScalarQuery) ExecuteScalarQuery(tickInfo runtime.TickInfo) (float64, error) {
	inPortReadings := runtime.PortToValue{}
	_, _ = scalarQuery.promQL.Execute(inPortReadings, tickInfo)
	// FYI: promQL ensures that initial runs return err when no queries have returned yet.
	return scalarQuery.promQL.value, scalarQuery.promQL.err
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
	componentIndex int,
	policyReadAPI iface.Policy,
	jobPostFix string,
) (*TaggedQuery, fx.Option, error) {
	scalarQuery, options, err := NewScalarQueryAndOptions(queryString, evaluationInterval, componentIndex, policyReadAPI, jobPostFix)
	if err != nil {
		return nil, fx.Options(), err
	}
	taggedQuery := &TaggedQuery{
		scalarQuery: scalarQuery,
		// Set err to make sure the initial runs of ExecutePromQuery return error.
		err: errNoQueriesReturned,
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
		promTimeout,
		taggedQuery.onTaggedResult,
	)
}

func (taggedQuery *TaggedQuery) onTaggedResult(res prometheusmodel.Value, err error) {
	taggedQuery.res = res
	taggedQuery.err = err
}

// ExecutePromQuery runs a PromQueryJob and returns the current results: res and err. This function is supposed to be run under Circuit Execution Lock (Execution of Circuit Components is protected by this lock).
func (taggedQuery *TaggedQuery) ExecutePromQuery(tickInfo runtime.TickInfo) (prometheusmodel.Value, error) {
	_, _ = taggedQuery.scalarQuery.ExecuteScalarQuery(tickInfo)
	return taggedQuery.res, taggedQuery.err
}
