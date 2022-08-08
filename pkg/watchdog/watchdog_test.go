package watchdog

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"testing"

	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/stretchr/testify/require"
)

const (
	// EnvTestIsolated is a marker property for the runner to confirm that this
	// test is running in isolation (i.e. a dedicated process).
	EnvTestIsolated = "TEST_ISOLATED"

	// EnvTestMemLimit is the memory limit applied when running test in container
	EnvTestMemLimit = "TEST_DOCKER_MEMLIMIT"
)

var (
	memLimit   int               // initialized in the init function
	limit64MiB uint64 = 64 << 20 // 64MiB.
	gws        jobs.GroupWatchers
	jws        jobs.JobWatchers
)

func init() {
	if l := os.Getenv(EnvTestMemLimit); l != "" {
		l, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		memLimit = l
	}
}

func skipIfNotIsolated(t *testing.T) {
	if os.Getenv(EnvTestIsolated) != "1" {
		t.Skipf("skipping test in non-isolated mode")
	}
}

// createJobGroup creates a job group and a multiJob
func createJobGroupAndMultiJob(key string, reg *status.Registry) (*jobs.JobGroup, *jobs.MultiJob) {
	group, err := jobs.NewJobGroup(key, reg, 0, jobs.RescheduleMode, gws)
	if err != nil {
		panic(fmt.Sprintf("Failed to create job group: %v", err))
	}
	multiJob := jobs.NewMultiJob(key+".MultiJob", "watchdogs", false, reg, jws, gws)
	group.Start()
	return group, multiJob
}

// runTest sets a GCPercent needed to trigger the enables policy and then overloads the memory, and ends up checking the memory stats to confirm that the test has run correctly.
func runTest(t *testing.T) {
	debug.SetGCPercent(100)
	rounds := 100
	if memLimit != 0 {
		rounds /= int(float64(memLimit)*0.8) / 1024 / 1024
	}

	// retain 1MiB every iteration.
	var retained [][]byte
	for i := 0; i < rounds; i++ {
		b := make([]byte, 1*1024*1024)
		for i := range b {
			b[i] = byte(i)
		}
		retained = append(retained, b)
	}

	for _, b := range retained {
		for i := range b {
			b[i] = byte(i)
		}
	}

	// Results are checked in a generic way.
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	require.NotZero(t, ms.NumGC)
	require.Zero(t, ms.NumForcedGC)
}

// TestIsolatedControlCGroup tests that the watchdog watermark and adaptive policy are enabled and working when the watchdog is run in an isolated process.
func TestIsolatedControlCGroup(t *testing.T) {
	skipIfNotIsolated(t)

	sentinel := newSentinel()
	var heapPolicy *heapPolicy
	key := "TestCGroupPolicy"
	reg := status.NewRegistry(".")

	jobGroup, multiJob := createJobGroupAndMultiJob(key, reg)
	config := WatchdogConfig{
		CGroup: WatchdogPolicyType{
			WatermarksPolicy: WatermarksPolicy{
				PolicyCommon: PolicyCommon{
					Enabled: true,
				},
			},
		},
	}
	watchDog := watchdog{
		statusRegistry: reg,
		jobGroup:       jobGroup,
		watchdogJob:    multiJob,
		config:         config,
	}
	err := setupWatchdogOnStart(context.Background(), watchDog, sentinel, heapPolicy)
	require.NoError(t, err)
	runTest(t)

	err = setupWatchdogOnStop(context.Background(), watchDog, sentinel)
	require.NoError(t, err)

	config.CGroup.WatermarksPolicy.PolicyCommon.Enabled = false
	config.CGroup.AdaptivePolicy.PolicyCommon.Enabled = true
	config.System.WatermarksPolicy.PolicyCommon.Enabled = false
	config.System.AdaptivePolicy.PolicyCommon.Enabled = true

	err = setupWatchdogOnStart(context.Background(), watchDog, sentinel, heapPolicy)
	require.NoError(t, err)
	runTest(t)

	err = setupWatchdogOnStop(context.Background(), watchDog, sentinel)
	require.NoError(t, err)
}

// TestIsolatedSystemDriven tests that the watchdog watermark and adaptive policy are enabled and working when the watchdog is run in an isolated process.
func TestIsolatedSystemDriven(t *testing.T) {
	skipIfNotIsolated(t)

	sentinel := newSentinel()
	var heapPolicy *heapPolicy
	key := "TestSystemDriven"
	reg := status.NewRegistry(".")
	jobGroup, multiJob := createJobGroupAndMultiJob(key, reg)

	config := WatchdogConfig{
		System: WatchdogPolicyType{
			AdaptivePolicy: AdaptivePolicy{
				PolicyCommon: PolicyCommon{
					Enabled: true,
				},
			},
		},
	}

	watchDog := watchdog{
		statusRegistry: reg,
		jobGroup:       jobGroup,
		watchdogJob:    multiJob,
		config:         config,
	}

	err := setupWatchdogOnStart(context.Background(), watchDog, sentinel, heapPolicy)
	require.NoError(t, err)
	runTest(t)

	err = setupWatchdogOnStop(context.Background(), watchDog, sentinel)
	require.NoError(t, err)

	config.System.AdaptivePolicy.PolicyCommon.Enabled = false
	config.System.WatermarksPolicy.PolicyCommon.Enabled = true

	err = setupWatchdogOnStart(context.Background(), watchDog, sentinel, heapPolicy)
	require.NoError(t, err)
	runTest(t)
}

// TestHeapDriven tests that the watchdog watermark and adaptive policy are enabled and working when the watchdog is run in an isolated process.
func TestIsolatedHeapDriven(t *testing.T) {
	skipIfNotIsolated(t)

	sentinel := newSentinel()
	var heapPolicy *heapPolicy
	key := "TestHeapPolicy"
	reg := status.NewRegistry(".")
	jobGroup, multiJob := createJobGroupAndMultiJob(key, reg)

	config := WatchdogConfig{
		Heap: HeapConfig{
			WatchdogPolicyType: WatchdogPolicyType{
				WatermarksPolicy: WatermarksPolicy{
					PolicyCommon: PolicyCommon{
						Enabled: true,
					},
				},
			},
			HeapLimit: HeapLimit{
				MinGoGC: 25,
				Limit:   268435456,
			},
		},
	}
	watchDog := watchdog{
		statusRegistry: reg,
		jobGroup:       jobGroup,
		watchdogJob:    multiJob,
		config:         config,
	}
	err := setupWatchdogOnStart(context.Background(), watchDog, sentinel, heapPolicy)
	require.NoError(t, err)
	runTest(t)

	err = setupWatchdogOnStop(context.Background(), watchDog, sentinel)
	require.NoError(t, err)
}
