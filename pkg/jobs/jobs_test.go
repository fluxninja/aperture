package jobs

import (
	"context"
	"time"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/status"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ JobWatcher = (*testGroupConfig)(nil)

type testJobRunConfig struct {
	expectedStatusMsg string        // expected Registry result for liveness check
	sleepTime         time.Duration // time to sleep before checking results
	expectedNoOfRuns  int           // Based on provided configuration, estimate the expected number of runs for the job
}

type testGroupConfig struct {
	jobs               []Job
	jobGroup           *JobGroup
	jobRunConfig       testJobRunConfig
	expectedScheduling bool // to be used to check if sleeping/stuck jobs are getting stuck or not
}

// When a job is only scheduled and not run, it's number of run should be 0
func (gc *testGroupConfig) OnJobScheduled() {
	for _, job := range gc.jobs {
		jobInfo := gc.jobGroup.JobInfo(job.Name())
		Expect(jobInfo.RunCount).To(Equal(0))
	}
}

// Estimate the number of runs for the job based on the provided configuration and compare with the actual result
func (gc *testGroupConfig) OnJobCompleted(_ *statusv1.Status, _ JobStats) {
	for _, job := range gc.jobs {
		jobInfo := gc.jobGroup.JobInfo(job.Name())
		Expect(jobInfo.RunCount).To(Equal(gc.jobRunConfig.expectedNoOfRuns))
	}
}

var _ = Describe("Jobs", func() {
	var (
		counter, counter2, counter3 int
		jobConfig                   JobConfig
		jobGroup                    *JobGroup
		groupWatchers               GroupWatchers
		jobwatchers                 JobWatchers
		addOne, addFive, addTen     func(ctx context.Context) (proto.Message, error)
	)

	BeforeEach(func() {
		var err error
		counter = 0
		counter2 = 0
		counter3 = 0

		alerter := alerts.NewSimpleAlerter(100)
		rootRegistry := status.NewRegistry(log.GetGlobalLogger(), alerter)
		jobConfig = JobConfig{
			InitialDelay:     config.MakeDuration(0),
			ExecutionPeriod:  config.MakeDuration(time.Millisecond * 200),
			ExecutionTimeout: config.MakeDuration(time.Millisecond * 200),
			InitiallyHealthy: false,
		}
		jobGroup, err = NewJobGroup(rootRegistry, 10, RescheduleMode, groupWatchers)
		Expect(err).To(BeNil())

		addOne = func(ctx context.Context) (proto.Message, error) {
			select {
			case <-ctx.Done():
				return &emptypb.Empty{}, nil
			default:
				counter += 1
				return &emptypb.Empty{}, nil
			}
		}

		addFive = func(ctx context.Context) (proto.Message, error) {
			select {
			case <-ctx.Done():
				return &emptypb.Empty{}, nil
			default:
				counter2 += 5
				return &emptypb.Empty{}, nil
			}
		}

		addTen = func(ctx context.Context) (proto.Message, error) {
			select {
			case <-ctx.Done():
				return &emptypb.Empty{}, nil
			default:
				counter3 += 10
				return &emptypb.Empty{}, nil
			}
		}

		// START JOB GROUP
		err = jobGroup.Start()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		// STOP JOB GROUP
		err := jobGroup.Stop()
		Expect(err).To(BeNil())
	})

	It("runs instant jobs", func() {
		job := NewBasicJob("instant-run-job", addOne)
		job.Execute(context.Background())
		Expect(counter).To(Equal(1))
	})

	It("runs instant jobs with specified number of runs", func() {
		job := NewBasicJob("instant-run-job", addOne)
		groupConfig := &testGroupConfig{
			jobs: []Job{job},
			jobRunConfig: testJobRunConfig{
				sleepTime:         time.Millisecond * 300,
				expectedStatusMsg: "OK",
				expectedNoOfRuns:  2,
			},
			expectedScheduling: true,
			jobGroup:           jobGroup,
		}
		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(2))
	})

	It("checks for timed out job", func() {
		jobConfig.InitialDelay = config.MakeDuration(time.Millisecond * 100)
		job := NewBasicJob("timeout-job",
			func(ctx context.Context) (proto.Message, error) {
				time.Sleep(time.Millisecond * 4000)
				select {
				case <-ctx.Done():
					return &emptypb.Empty{}, nil
				default:
					counter += 1
					return &emptypb.Empty{}, nil
				}
			})

		groupConfig := &testGroupConfig{
			jobs: []Job{job},
			jobRunConfig: testJobRunConfig{
				sleepTime:         time.Millisecond * 2000,
				expectedStatusMsg: "Timeout",
				expectedNoOfRuns:  11, // Job gets scheduled 10 times within 2 seconds + 1 scheduling caused by manual trigger if job is stuck
			},
			expectedScheduling: false,
			jobGroup:           jobGroup,
		}
		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(0))
	})

	It("runs multiple basic jobs", func() {
		job := NewBasicJob("test-job", addOne)
		job2 := NewBasicJob("test-job2", addFive)

		groupConfig := &testGroupConfig{
			jobs: []Job{job, job2},
			jobRunConfig: testJobRunConfig{
				sleepTime:         time.Millisecond * 300,
				expectedStatusMsg: "OK",
				expectedNoOfRuns:  2,
			},
			expectedScheduling: true,
			jobGroup:           jobGroup,
		}
		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(2))
		Expect(counter2).To(Equal(10))
	})

	It("runs multi jobs", func() {
		multiJob := NewMultiJob(jobGroup.GetStatusRegistry().Child("multi-job", "multi-job"), jobwatchers, groupWatchers)
		job := NewBasicJob("test-job", addOne)
		job2 := NewBasicJob("test-job2", addFive)

		groupConfig := &testGroupConfig{
			jobs: []Job{multiJob},
			jobRunConfig: testJobRunConfig{
				sleepTime:         time.Millisecond * 300,
				expectedStatusMsg: "OK",
				expectedNoOfRuns:  2,
			},
			expectedScheduling: true,
			jobGroup:           jobGroup,
		}
		multiJob.RegisterJob(job)
		multiJob.RegisterJob(job2)

		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(2))
		Expect(counter2).To(Equal(10))

		multiJob.DeregisterAll()
	})

	It("runs multiple multi jobs", func() {
		multiJob := NewMultiJob(jobGroup.GetStatusRegistry().Child("multi-job1", "multi-job1"), jobwatchers, groupWatchers)
		multiJob2 := NewMultiJob(jobGroup.GetStatusRegistry().Child("multi-job2", "multi-job2"), jobwatchers, groupWatchers)
		job := NewBasicJob("test-job", addOne)
		job2 := NewBasicJob("test-job2", addFive)
		job3 := NewBasicJob("test-job3", addTen)

		groupConfig := &testGroupConfig{
			jobs: []Job{multiJob, multiJob2},
			jobRunConfig: testJobRunConfig{
				sleepTime:         time.Millisecond * 300,
				expectedStatusMsg: "OK",
				expectedNoOfRuns:  2,
			},
			expectedScheduling: true,
			jobGroup:           jobGroup,
		}
		multiJob.RegisterJob(job)
		multiJob.RegisterJob(job2)
		multiJob2.RegisterJob(job3)

		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(2))
		Expect(counter2).To(Equal(10))
		Expect(counter3).To(Equal(20))

		multiJob.DeregisterAll()
		multiJob2.DeregisterAll()
	})

	It("does not allow running nil jobs", func() {
		job := NewBasicJob("instant-run-job", nil)
		_, err := job.Execute(context.Background())
		Expect(err).ToNot(BeNil())
	})

	It("allows registering nil jobs", func() {
		job := NewBasicJob("test-job", nil)
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).To(BeNil())
	})

	It("does not allow registering same job", func() {
		job := NewBasicJob("test-job", addOne)
		job2 := job
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).To(BeNil())

		err = jobGroup.RegisterJob(job2, jobConfig)
		Expect(err).ToNot(BeNil())

		jobGroup.DeregisterAll()
		// error when registering job multiple times, written here to achieve more coverage
		err = jobGroup.DeregisterJob(job.Name())
		Expect(err).ToNot(BeNil())
		Expect(jobGroup.JobInfo(job.Name())).To(BeNil())
	})
})

func runJobGroup(jobGroup *JobGroup, groupConfig *testGroupConfig, jobConfig JobConfig) {
	// Register all jobs in group
	for _, job := range groupConfig.jobs {
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).To(BeNil())
	}

	// check initially healthy status
	Expect(jobGroup.IsHealthy()).To(Equal(jobConfig.InitiallyHealthy))
	time.Sleep(groupConfig.jobRunConfig.sleepTime)

	for _, job := range groupConfig.jobs {
		registry := jobGroup.livenessRegistry
		livenessReg := registry.Root().
			Child("subsystem", "liveness").
			Child("jg", "job_groups").
			Child(registry.Key(), registry.Value()).
			Child("executor", job.Name())

		if groupConfig.jobRunConfig.expectedStatusMsg == "Timeout" {
			checkStatusMessage(livenessReg, groupConfig.jobRunConfig.expectedStatusMsg, true)
			jobGroup.TriggerJob(job.Name())
		} else {
			checkStatusMessage(livenessReg, groupConfig.jobRunConfig.expectedStatusMsg, false)
		}

		groupConfig.OnJobCompleted(nil, JobStats{})
		_, val := jobGroup.Results()
		Expect(val).To(Equal(groupConfig.expectedScheduling))
	}

	// cleanup
	for _, job := range groupConfig.jobs {
		err := jobGroup.DeregisterJob(job.Name())
		Expect(err).To(BeNil())
	}
}

func checkStatusMessage(registry status.Registry, protoMsg string, hasError bool) {
	var gotStatusMsg, expectedStatusMsg *anypb.Any

	Expect(registry.HasError()).To(Equal(hasError))

	gotStatusMsg = registry.GetStatus().GetMessage()

	expectedStatusMsg, _ = anypb.New(&emptypb.Empty{})
	if protoMsg != "" {
		expectedStatusMsg, _ = anypb.New(wrapperspb.String(protoMsg))
	}
	Expect(proto.Equal(gotStatusMsg, expectedStatusMsg)).To(BeTrue())
}
