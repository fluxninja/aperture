package harness

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	// EtcdBinPath is the path to the etcd binary.
	EtcdBinPath      = "etcd"
	etcdLocalAddress = "127.0.0.1:0"
)

// EtcdHarness represents a running etcd server for an integration test environment.
type EtcdHarness struct {
	errWriter  io.Writer
	etcdServer *exec.Cmd
	etcdDir    string
	Client     *clientv3.Client
	Endpoint   string
}

// NewEtcdHarness initializes a harnessed etcd server and returns the EtcdHarness.
func NewEtcdHarness(etcdErrWriter io.Writer) (*EtcdHarness, error) {
	h := &EtcdHarness{
		errWriter: etcdErrWriter,
	}

	endpointAddr, err := AllocateLocalAddress(etcdLocalAddress)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("http://%s", endpointAddr)

	peerAddr, err := AllocateLocalAddress(etcdLocalAddress)
	if err != nil {
		return nil, err
	}
	peer := fmt.Sprintf("http://%s", peerAddr)

	etcdBin, err := LocalBinAvailable(EtcdBinPath)
	if err != nil {
		return nil, err
	}

	h.etcdDir, err = os.MkdirTemp("/tmp", "etcd_testserver")
	if err != nil {
		return nil, err
	}

	h.etcdServer = exec.Command(
		etcdBin,
		"--data-dir="+h.etcdDir,
		"--listen-peer-urls="+peer,
		"--initial-cluster="+"default="+peer,
		"--initial-advertise-peer-urls="+peer,
		"--listen-client-urls="+endpoint,
		"--advertise-client-urls="+endpoint,
	)
	h.etcdServer.Stderr = h.errWriter
	h.etcdServer.Stdout = io.Discard
	h.Endpoint = endpoint

	err = h.etcdServer.Start()
	if err != nil {
		h.Stop()
		return nil, err
	}

	h.Client, err = clientv3.New(clientv3.Config{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		h.Stop()
		return h, err
	}

	err = h.pollEtcdForReadiness()
	if err != nil {
		h.Stop()
		return nil, err
	}

	return h, nil
}

func (h *EtcdHarness) pollEtcdForReadiness() error {
	// Actively poll for etcd coming up for 4 seconds every 200 milliseconds.
	for i := 0; i < 20; i++ {
		until := time.Now().Add(200 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.TODO(), until)
		_, err := h.Client.Get(ctx, "/")
		cancel()
		if err == nil {
			return nil
		}
		toSleep := time.Until(until)
		if toSleep > 0 {
			time.Sleep(toSleep)
		}
	}
	return fmt.Errorf("etcd didn't come up in 4000ms")
}

// Stop kills the harnessed etcd server and cleans up the etcd directory.
func (h *EtcdHarness) Stop() {
	if h.etcdServer != nil {
		_ = h.etcdServer.Process.Kill()
		_ = h.etcdServer.Wait()
	}
	if h.etcdDir != "" {
		_ = os.RemoveAll(h.etcdDir)
	}
}
