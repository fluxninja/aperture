//go:build !linux
// +build !linux

package watchdog

import (
	"context"

	"google.golang.org/protobuf/proto"

	watchdogconfig "github.com/fluxninja/aperture/pkg/watchdog/config"
)

type cgroupWatermarks struct {
	watchdogconfig.WatermarksPolicy
}

type cgroupAdaptive struct {
	watchdogconfig.AdaptivePolicy
}

// Valid checks if the policy is valid in linux. Therefore, it always returns false in non-linux builds.
func (policy *cgroupWatermarks) Valid() bool {
	return false
}

// Check evaluates the cgroup memory usage and runs GC at configured watermarks of memory utilization.
// This will return nil, nil in non-linux builds.
func (policy *cgroupWatermarks) Check(ctx context.Context) (proto.Message, error) {
	return nil, nil
}

// Valid checks if the policy is valid in linux. Therefore, it always returns false in non-linux builds.
func (policy *cgroupAdaptive) Valid() bool {
	return false
}

// Check evaluates the cgroup memory usage and runs GC at configured adaptive thresholds of memory utilization.
// This will return nil, nil in non-linux builds.
func (policy *cgroupAdaptive) Check(ctx context.Context) (proto.Message, error) {
	return nil, nil
}
