package distcache

import (
	"context"
	"fmt"
	stdlog "log"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/hashicorp/memberlist"
)

func newTestDistCacheWithConfig(t *testing.T, c *olricconfig.Config) (*DistCache, error) {
	distCache := &DistCache{}

	ctx, cancel := context.WithCancel(context.Background())

	if c != nil {
		distCache.config = c
	} else {
		oc := olricconfig.New("local")

		err := oc.DMaps.Sanitize()
		if err != nil {
			t.Errorf("Failed to sanitize olric config: %v", err)
		}
		err = oc.DMaps.Validate()
		if err != nil {
			t.Errorf("Failed to validate olric config: %v", err)
		}
		distCache.config = oc
	}

	distCache.config.Started = func() {
		t.Log("Started olric server")
		cancel()
	}

	o, err := olric.New(distCache.config)
	if err != nil {
		return nil, err
	}

	distCache.olric = o

	go func() {
		t.Log("Starting DistCacheLimiter")
		err = distCache.olric.Start()
		if err != nil {
			t.Errorf("Failed to start olric: %v", err)
		}
	}()

	select {
	case <-time.After(time.Second):
		t.Fatal("DistCache cannot be started in one second")
	case <-ctx.Done():
		// everything is fine
	}

	return distCache, nil
}

// TestDistCacheCluster is a test cluster of DistCache instances.
type TestDistCacheCluster struct {
	Lock    sync.Mutex
	Members map[string]*DistCache
}

func newTestOlricConfig() *olricconfig.Config {
	c := olricconfig.New("local")
	mc := memberlist.DefaultLocalConfig()
	mc.BindAddr = "127.0.0.1"
	mc.BindPort = 0
	c.MemberlistConfig = mc
	c.Logger = stdlog.New(&OlricLogWriter{Logger: log.GetGlobalLogger()}, "", 0)

	port, err := getFreePort()
	if err != nil {
		panic(fmt.Sprintf("GetFreePort returned an error: %v", err))
	}
	c.BindAddr = "127.0.0.1"
	c.BindPort = port
	c.MemberlistConfig.Name = net.JoinHostPort(c.BindAddr, strconv.Itoa(c.BindPort))
	if err := c.Sanitize(); err != nil {
		panic(fmt.Sprintf("failed to sanitize default config: %v", err))
	}
	return c
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func (cl *TestDistCacheCluster) addDistCacheWithConfig(t *testing.T, c *olricconfig.Config) error {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()

	if c == nil {
		c = newTestOlricConfig()
	}

	for _, member := range cl.Members {
		c.Peers = append(c.Peers, fmt.Sprintf("%s:%d", member.config.MemberlistConfig.BindAddr, member.config.MemberlistConfig.BindPort))
	}

	dc, err := newTestDistCacheWithConfig(t, c)
	if err != nil {
		return err
	}
	cl.Members[dc.config.MemberlistConfig.Name] = dc
	return nil
}

// NewTestDistCacheCluster creates a new DistCache cluster with n members.
func NewTestDistCacheCluster(t *testing.T, n int) *TestDistCacheCluster {
	cl := &TestDistCacheCluster{
		Members: make(map[string]*DistCache, n),
	}
	for i := 0; i < n; i++ {
		t.Log("Adding cluster member")
		_ = cl.addDistCacheWithConfig(t, nil)
	}
	return cl
}

// CloseTestDistCacheCluster closes the test dist cache cluster.
func CloseTestDistCacheCluster(t *testing.T, cl *TestDistCacheCluster) {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()
	for _, member := range cl.Members {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := member.olric.Shutdown(ctx)
		if err != nil {
			t.Log("Failed to shutdown olric")
		}
	}
}
