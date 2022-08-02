package watchdog

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/fluxninja/aperture/pkg/config"
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

// func main() {
// 	app := platform.New(
// 		platform.Config{}.Module(),
// 		http.ClientConstructor{Name: "k8s-http-client", Key: "kubernetes.http_client"}.Annotate(),
// 		notifiers.TrackersConstructor{Name: "entity_trackers"}.Annotate(),
// 		prometheus.Module(),
// 		otel.ProvideAnnotatedAgentConfig(),
// 		fx.Provide(
// 			agentinfo.ProvideAgentInfo,
// 			k8s.Providek8sClient,
// 			clockwork.NewRealClock,
// 			entitycache.ProvideEntityCache,
// 			otel.AgentOTELComponents,
// 			agent.ProvidePeersPrefix,
// 		),
// 		flowcontrol.Module,
// 		classification.Module,
// 		authz.Module,
// 		otelcollector.Module(),
// 		distcache.Module(),
// 		dataplane.PolicyModule(),
// 		discovery.Module(),
// 		fx.Invoke(
// 			authz.Register,
// 			flowcontrol.Register,
// 		),
// 		grpc.ClientConstructor{Name: "flowcontrol-grpc-client", Key: "flowcontrol.client.grpc"}.Annotate(),
// 	)

// 	if err := app.Err(); err != nil {
// 		visualize, _ := fx.VisualizeError(err)
// 		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
// 	}

// 	log.Info().Msg("aperture-agent app created")
// 	platform.Run(app)
// }

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

func createUnmarshaller() (config.Unmarshaller, error) {
	unmarshaller, err := config.KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller([]byte(""))
	return unmarshaller, err
}

func TestIsolatedControlWatermarkPolicy(t *testing.T) {
	key := "test-isolated-control-watermark-policy"
	reg := status.NewRegistry(".")
	// skipIfNotIsolated(t)

	jobGroup, multiJob := createJobGroupAndMultiJob(key, reg)
	unmarshaller, err := createUnmarshaller()
	require.NoError(t, err)

	constructor := &Constructor{
		Key: key,
		DefaultConfig: WatchdogConfig{
			CGroup: WatchdogPolicyType{
				WatermarksPolicy: WatermarksPolicy{
					Watermarks: []float64{0.50, 0.75, 0.80, 0.85, 0.90, 0.95, 0.99},
					PolicyCommon: PolicyCommon{
						Enabled: true,
					},
				},
			},
		},
	}

	watchDog := WatchdogIn{
		StatusRegistry: status.NewRegistry("."),
		JobGroup:       jobGroup,
		WatchdogJob:    multiJob,
		Unmarshaller:   unmarshaller,
	}

	err = constructor.setupWatchdog(watchDog)
	require.NoError(t, err)

	fmt.Printf("watchdog Content: %v\n", watchDog)

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

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	require.NotZero(t, ms.NumGC)    // GCs have taken place, but...
	require.Zero(t, ms.NumForcedGC) // ... no forced GCs beyond our initial one.
}
