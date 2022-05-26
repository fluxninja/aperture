package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"

	"aperture.tech/aperture/pkg/agentinfo"
	"aperture.tech/aperture/pkg/authz"
	"aperture.tech/aperture/pkg/classification"
	"aperture.tech/aperture/pkg/entitycache"
	etcdclient "aperture.tech/aperture/pkg/etcd/client"
	etcdwatcher "aperture.tech/aperture/pkg/etcd/watcher"
	"aperture.tech/aperture/pkg/flowcontrol"
	"aperture.tech/aperture/pkg/jobs"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/net/grpc"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/otel"
	"aperture.tech/aperture/pkg/otelcollector"
	"aperture.tech/aperture/pkg/platform"
	"aperture.tech/aperture/pkg/policies/dataplane"
	"aperture.tech/aperture/pkg/status"
	"aperture.tech/aperture/pkg/utils"
	"aperture.tech/aperture/test/harness"
)

const (
	jobGroupName = "job-group"
)

var (
	project     string
	app         *fx.App
	addr        string
	configDir   string
	l           *utils.GoLeakDetector
	eh          *harness.EtcdHarness
	ehStarted   bool
	etcdClient  *etcdclient.Client
	etcdWatcher notifiers.Watcher
	ph          *harness.PrometheusHarness
	phStarted   bool
	jgIn        *JobGroupIn
)

type JobGroupIn struct {
	fx.In
	JobGroup *jobs.JobGroup `name:"job-group"`
	Registry *status.Registry
}

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aperture Suite")
}

var _ = BeforeSuite(func() {
	var err error

	addr = ":18081"

	configDir = "aperturetest-config"

	l = utils.NewGoLeakDetector()

	_, err = harness.LocalBinAvailable(harness.EtcdBinPath)
	if err == nil {
		eh, err = harness.NewEtcdHarness(os.Stderr)
		Expect(err).NotTo(HaveOccurred())
		ehStarted = true
	}

	_, err = harness.LocalBinAvailable(harness.PrometheusBinPath)
	if err == nil {
		ph, err = harness.NewPrometheusHarness(os.Stderr)
		Expect(err).NotTo(HaveOccurred())
		phStarted = true
	}

	// reset config dir
	err = os.RemoveAll(configDir)
	Expect(err).NotTo(HaveOccurred())
	err = os.MkdirAll(configDir, 0o777)
	Expect(err).NotTo(HaveOccurred())

	jgIn = &JobGroupIn{}

	apertureConfig := map[string]interface{}{
		"plugins": map[string]interface{}{
			"disabled_plugins": []string{"aperture-plugin-fluxninja"},
		},
		"log": map[string]interface{}{
			// for cleaner logs and for testing config
			"level":          "debug",
			"pretty_console": true,
		},
		"server": map[string]interface{}{
			"addr": addr,
		},
		"config_path": configDir,
		"sentrywriter": map[string]interface{}{
			"disabled": true,
		},
	}

	if ehStarted {
		apertureConfig["etcd"] = map[string]interface{}{
			"endpoints": []string{eh.Endpoint},
		}
	}

	if phStarted {
		apertureConfig["prometheus"] = map[string]interface{}{
			"address": ph.Endpoint,
		}
	}

	apertureOpts := fx.Options(
		platform.Config{MergeConfig: apertureConfig}.Module(),
		fx.Option(
			fx.Provide(
				fx.Annotate(
					provideOTELConfig,
					fx.ResultTags(`name:"base"`),
				),
			),
		),
		fx.Provide(
			clockwork.NewRealClock,
			otel.AgentOTELComponents,
			dataplane.ProvideEngineAPI,
			entitycache.NewEntityCache,
			agentinfo.ProvideAgentInfo,
		),
		flowcontrol.Module,
		classification.Module,
		authz.Module,
		otelcollector.Module(),
		fx.Invoke(
			authz.Register,
			flowcontrol.Register,
		),
		grpc.ClientConstructor{Name: "flowcontrol-grpc-client", Key: "flowcontrol.client.grpc"}.Annotate(),
		jobs.JobGroupConstructor{Group: jobGroupName}.Annotate(),
		fx.Populate(jgIn),
	)

	if ehStarted {
		apertureOpts = fx.Options(
			apertureOpts,
			etcdwatcher.Constructor{Name: "test-etcd-watcher", EtcdPath: "foo/"}.Annotate(),
			fx.Populate(&etcdClient),
		)
	}

	app = platform.New(apertureOpts)

	err = app.Err()

	if err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Error().Err(err).Msg("fx.New failed: " + visualize)
	}

	Expect(err).NotTo(HaveOccurred())

	startCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = app.Start(startCtx)
	Expect(err).NotTo(HaveOccurred())

	etcdWatcher, err = etcdwatcher.NewWatcher(etcdClient, "foo/")
	Expect(err).NotTo(HaveOccurred())
	etcdWatcher.Start()

	project = "staging"
	Eventually(func() bool {
		_, err := http.Get(fmt.Sprintf("http://%v/version", addr))
		return err == nil
	}).Should(BeTrue())
})

var _ = AfterSuite(func() {
	stopCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error

	if etcdWatcher != nil {
		etcdWatcher.Stop()
	}

	err = app.Stop(stopCtx)
	Expect(err).NotTo(HaveOccurred())

	_ = os.RemoveAll(configDir)

	if ehStarted {
		eh.Stop()
	}

	if phStarted {
		ph.Stop()
	}

	err = l.FindLeaks()
	Expect(err).NotTo(HaveOccurred())
})

func provideOTELConfig() *otelcollector.OTELConfig {
	cfg := otelcollector.NewOTELConfig()
	cfg.AddReceiver("prometheus", map[string]interface{}{
		"config": map[string]interface{}{
			"scrape_configs": []map[string]interface{}{
				{
					"job_name":        "aperture-agent",
					"scrape_interval": "5s",
					"static_configs": []map[string]interface{}{
						{
							"targets": []string{addr},
						},
					},
				},
			},
		},
	})
	cfg.AddExporter("prometheusremotewrite", map[string]interface{}{
		"endpoint": ph.Endpoint + "/api/v1/write",
	})
	cfg.Service.AddPipeline("metrics", otelcollector.Pipeline{
		Receivers: []string{"prometheus"},
		Exporters: []string{"prometheusremotewrite"},
	})
	return cfg
}
