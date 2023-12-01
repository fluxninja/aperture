package jobs

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/status"
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

// When a job is only scheduled and not run, it is number of run should be 0
func (gc *testGroupConfig) OnJobScheduled() {
	for _, job := range gc.jobs {
		jobInfo, err := gc.jobGroup.JobInfo(job.Name())
		Expect(err).NotTo(HaveOccurred())
		Expect(jobInfo.ExecuteCount).To(Equal(0))
	}
}

// Estimate the number of runs for the job based on the provided configuration and compare with the actual result
func (gc *testGroupConfig) OnJobCompleted(_ *statusv1.Status, _ JobStats) {
	for _, job := range gc.jobs {
		jobInfo, err := gc.jobGroup.JobInfo(job.Name())
		Expect(err).NotTo(HaveOccurred())
		Expect(jobInfo.ExecuteCount).To(Equal(gc.jobRunConfig.expectedNoOfRuns))
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
			ExecutionPeriod:  config.MakeDuration(time.Millisecond * 200),
			ExecutionTimeout: config.MakeDuration(time.Millisecond * 200),
			InitiallyHealthy: false,
		}
		jobGroup, err = NewJobGroup(rootRegistry, JobGroupConfig{}, groupWatchers)
		Expect(err).NotTo(HaveOccurred())

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
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// STOP JOB GROUP
		err := jobGroup.Stop()
		Expect(err).NotTo(HaveOccurred())
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
				sleepTime:         time.Millisecond * 1000,
				expectedStatusMsg: "OK",
				expectedNoOfRuns:  5,
			},
			expectedScheduling: true,
			jobGroup:           jobGroup,
		}
		runJobGroup(jobGroup, groupConfig, jobConfig)
		Expect(counter).To(Equal(5))
	})

	It("checks for timed out job", func() {
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
				expectedNoOfRuns:  1,
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
				sleepTime:         time.Millisecond * 500,
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
				sleepTime:         time.Millisecond * 500,
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
				sleepTime:         time.Millisecond * 500,
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
		Expect(err).To(HaveOccurred())
	})

	It("allows registering nil jobs", func() {
		job := NewBasicJob("test-job", nil)
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("does not allow registering same job", func() {
		job := NewBasicJob("test-job", addOne)
		job2 := job
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).NotTo(HaveOccurred())

		err = jobGroup.RegisterJob(job2, jobConfig)
		Expect(err).To(HaveOccurred())

		jobGroup.DeregisterAll()
		// error when registering job multiple times, written here to achieve more coverage
		err = jobGroup.DeregisterJob(job.Name())
		Expect(err).To(HaveOccurred())
		_, err = jobGroup.JobInfo(job.Name())
		Expect(err).To(HaveOccurred())
	})
})

func runJobGroup(jobGroup *JobGroup, groupConfig *testGroupConfig, jobConfig JobConfig) {
	// Register all jobs in group
	for _, job := range groupConfig.jobs {
		err := jobGroup.RegisterJob(job, jobConfig)
		Expect(err).NotTo(HaveOccurred())
	}

	// check initially healthy status
	Expect(jobGroup.IsHealthy()).To(Equal(jobConfig.InitiallyHealthy))
	time.Sleep(groupConfig.jobRunConfig.sleepTime)

	for _, job := range groupConfig.jobs {
		registry := jobGroup.livenessRegistry
		livenessReg := registry.Root().
			Child("system", "liveness").
			Child("job-group", "job_groups").
			Child(registry.Key(), registry.Value()).
			Child("executor", job.Name())

		if groupConfig.jobRunConfig.expectedStatusMsg == "Timeout" {
			checkStatusMessage(livenessReg, groupConfig.jobRunConfig.expectedStatusMsg, true)
			jobGroup.TriggerJob(job.Name(), time.Duration(0))
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
		Expect(err).NotTo(HaveOccurred())
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
