package harness

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	// PrometheusBinPath is the path to the prometheus binary.
	PrometheusBinPath = "prometheus"
	// PrometheusLocalAddress is the local address to bind prometheus to.
	PrometheusLocalAddress = "127.0.0.1:9095"
	prometheusConfig       = "./harness/prometheus.yml"
)

// PrometheusHarness represents a running prometheus server for an integration test environment.
type PrometheusHarness struct {
	errWriter        io.Writer
	prometheusServer *exec.Cmd
	prometheusDir    string
	client           prometheusv1.API
	Endpoint         string
}

// NewPrometheusHarness initializes a harnessed prometheus server and returns the PrometheusHarness.
func NewPrometheusHarness(prometheusWriter io.Writer) (*PrometheusHarness, error) {
	p := &PrometheusHarness{
		errWriter: prometheusWriter,
	}

	endpointAddr, err := AllocateLocalAddress(PrometheusLocalAddress)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("http://%s", endpointAddr)

	prometheusBin, err := LocalBinAvailable(PrometheusBinPath)
	if err != nil {
		return nil, err
	}

	p.prometheusDir, err = os.MkdirTemp("/tmp", "prometheus_testserver")
	if err != nil {
		return nil, err
	}

	// Bring up prometheus:
	// prometheus --config.file="./harness/prometheus.yml" --web.route-prefix="/" --web.listen-address="127.0.0.1:9095"
	// --web.enable-remote-write-receiver --web.enable-lifecycle --storage.tsdb.path="/tmp/prometheus_testserver"
	p.prometheusServer = exec.Command(
		prometheusBin,
		"--config.file="+prometheusConfig,
		"--web.route-prefix="+"/",
		"--web.listen-address="+endpointAddr,
		"--web.enable-remote-write-receiver",
		"--web.enable-lifecycle",
		"--storage.tsdb.path="+p.prometheusDir,
	)
	p.prometheusServer.Stderr = p.errWriter
	p.prometheusServer.Stdout = io.Discard
	p.Endpoint = endpoint

	err = p.prometheusServer.Start()
	if err != nil {
		p.Stop()
		return nil, err
	}

	var apiClient api.Client
	apiClient, err = api.NewClient(api.Config{
		Address: endpoint,
	})
	if err != nil {
		p.Stop()
		return nil, err
	}
	p.client = prometheusv1.NewAPI(apiClient)

	err = p.pollPrometheusForReadiness()
	if err != nil {
		p.Stop()
		return nil, err
	}

	return p, nil
}

func (p *PrometheusHarness) pollPrometheusForReadiness() error {
	// Actively poll for prometheus coming up for 4 seconds every 200 milliseconds.
	for i := 0; i < 20; i++ {
		until := time.Now().Add(200 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.TODO(), until)
		_, _, err := p.client.Query(ctx, "up", time.Now())
		cancel()
		if err == nil {
			return nil
		}
		toSleep := time.Until(until)
		if toSleep > 0 {
			time.Sleep(toSleep)
		}
	}
	return fmt.Errorf("prometheus didn't come up in 4000ms")
}

// Stop kills the harnessed prometheus server and cleans up the prometheus directory.
func (p *PrometheusHarness) Stop() {
	if p.prometheusServer != nil {
		_ = p.prometheusServer.Process.Kill()
		_ = p.prometheusServer.Wait()
	}
	if p.prometheusDir != "" {
		_ = os.RemoveAll(p.prometheusDir)
	}
}
