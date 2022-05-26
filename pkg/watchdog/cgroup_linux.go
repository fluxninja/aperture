package watchdog

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/cgroups"
	"google.golang.org/protobuf/proto"

	"aperture.tech/aperture/pkg/log"
)

type cgroupBase struct {
	cgroup cgroups.Cgroup
}

func (cg *cgroupBase) load() error {
	if cg.cgroup == nil {
		// use self path unless our PID is 1, in which case we're running inside
		// a container and our limits are in the root path.
		path := cgroups.NestedPath("")
		if pid := os.Getpid(); pid == 1 {
			path = cgroups.RootPath
		}

		cgroup, err := cgroups.Load(cgroups.SingleSubsystem(cgroups.V1, cgroups.Memory), path)
		if err != nil {
			return fmt.Errorf("failed to load cgroup for process: %w", err)
		}
		cg.cgroup = cgroup
	}
	return nil
}

func (cg *cgroupBase) cgroupUsage() (uint64, uint64, error) {
	if err := cg.load(); err != nil {
		return 0, 0, err
	}

	stat, err := cg.cgroup.Stat()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to load memory cgroup stats: %w", err)
	} else if stat.Memory == nil || stat.Memory.Usage == nil {
		return 0, 0, fmt.Errorf("cgroup memory stats are nil; aborting")
	}
	return stat.Memory.Usage.Limit, stat.Memory.Usage.Usage, nil
}

type cgroupWatermarks struct {
	cgroupBase
	WatermarksPolicy
}

type cgroupAdaptive struct {
	cgroupBase
	AdaptivePolicy
}

// Valid checks if the policy is valid in linux.
func (policy *cgroupWatermarks) Valid() bool {
	return true
}

// Check evaluates the cgroup memory usage and runs GC at configured watermarks of memory utilization.
func (policy *cgroupWatermarks) Check(ctx context.Context) (proto.Message, error) {
	log.Debug().Msg("CGroup watermarks check triggered")
	return check(policy, ctx, policy.cgroupUsage)
}

// Valid checks if the policy is valid in linux.
func (policy *cgroupAdaptive) Valid() bool {
	return true
}

// Check evaluates the cgroup memory usage and runs GC at configured adaptive thresholds of memory utilization.
func (policy *cgroupAdaptive) Check(ctx context.Context) (proto.Message, error) {
	log.Debug().Msg("CGroup watermarks check triggered")
	return check(policy, ctx, policy.cgroupUsage)
}
