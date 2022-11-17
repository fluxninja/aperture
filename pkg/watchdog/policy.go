// +kubebuilder:validation:Optional
package watchdog

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	watchdogv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/watchdog/v1"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
)

// WatchdogConfig holds configuration for Watchdog Policy. For each policy, either watermark or adaptive should be configured.
// swagger:model
// +kubebuilder:object:generate=true
type WatchdogConfig struct {
	Job jobs.JobConfig `json:"job"`

	CGroup WatchdogPolicyType `json:"cgroup"`

	System WatchdogPolicyType `json:"system"`

	Heap HeapConfig `json:"heap"`
}

// WatchdogPolicyType holds configuration Watchdog Policy algorithms. If both algorithms are configured then only watermark algorithm is used.
// swagger:model WatchdogPolicyType
// +kubebuilder:object:generate=true
type WatchdogPolicyType struct {
	WatermarksPolicy WatermarksPolicy `json:"watermarks_policy"`

	AdaptivePolicy AdaptivePolicy `json:"adaptive_policy"`
}

// HeapLimit holds configuration for Watchdog heap limit.
// swagger:model
// +kubebuilder:object:generate=true
type HeapLimit struct {
	// Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid overscheduling.
	MinGoGC int `json:"min_gogc" validate:"gt=0,lte=100" default:"25"`

	// Maximum memory (in bytes) sets limit of process usage. Default = 256MB.
	Limit uint64 `json:"limit" validate:"gt=0" default:"268435456"`
}

// HeapConfig holds configuration for heap Watchdog.
// swagger:model
// +kubebuilder:object:generate=true
type HeapConfig struct {
	WatchdogPolicyType `json:",inline"`

	HeapLimit `json:",inline"`
}

// PolicyCommon holds common configuration for Watchdog policies.
// swagger:model
// +kubebuilder:object:generate=true
type PolicyCommon struct {
	// Flag to enable the policy
	Enabled bool `json:"enabled" default:"false"`

	// total
	// swagger:ignore
	total uint64
}

type usageFn func() (total uint64, usage uint64, err error)

type policyInterface interface {
	nextThreshold(total, used uint64) uint64
}

// WatermarksPolicy creates a Watchdog policy that schedules GC at concrete watermarks.
// swagger:model
// +kubebuilder:object:generate=true
type WatermarksPolicy struct {
	// Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It is recommended to set an extreme watermark for the last element (e.g. 0.99).
	Watermarks []float64 `json:"watermarks,omitempty" validate:"omitempty,dive,gte=0,lte=1" default:"[0.50,0.75,0.80,0.85,0.90,0.95,0.99]"`

	// internal fields
	// swagger:ignore
	thresholds []uint64

	PolicyCommon `json:",inline"`
}

func (policy *WatermarksPolicy) nextThreshold(total, used uint64) uint64 {
	if len(policy.thresholds) == 0 || policy.total != total {
		policy.total = total
		// sort watermarks
		sort.Float64s(policy.Watermarks)
		policy.thresholds = make([]uint64, 0, len(policy.Watermarks))
		for _, m := range policy.Watermarks {
			threshold := uint64(float64(policy.total) * m)
			policy.thresholds = append(policy.thresholds, threshold)
		}
		log.Info().Floats64("watermarks", policy.Watermarks).Uints64("thresholds", policy.thresholds).Msg("Initialized watermark watchdog policy")
	}
	var i int
	for ; i < len(policy.thresholds); i++ {
		t := policy.thresholds[i]
		if used < t {
			return t
		}
	}
	// we reached the maximum threshold, so we disable this policy.
	return policyTempDisabled
}

// AdaptivePolicy creates a policy that forces GC when the usage surpasses the configured factor of the available memory. This policy calculates next target as usage+(limit-usage)*factor.
// swagger:model
// +kubebuilder:object:generate=true
type AdaptivePolicy struct {
	PolicyCommon `json:",inline"`

	// Factor sets user-configured limit of available memory
	Factor float64 `json:"factor" validate:"gte=0,lte=1" default:"0.50"`
}

func (policy *AdaptivePolicy) nextThreshold(total, used uint64) uint64 {
	if policy.total != total {
		policy.total = total
	}

	if used > policy.total {
		return used
	}
	available := float64(policy.total) - float64(used)
	next := used + uint64(available*policy.Factor)
	return next
}

func forceGC() time.Duration {
	log.Info().Msg("watchdog is forcing GC")
	start := time.Now()
	runtime.GC()
	took := time.Since(start)
	log.Info().Dur("took", took).Msg("Watchdog triggered GC finished")
	return took
}

func check(policy policyInterface, _ context.Context, fn usageFn) (proto.Message, error) {
	total, used, err := fn()
	if err != nil {
		return nil, err
	}

	threshold := policy.nextThreshold(total, used)

	result := &watchdogv1.WatchdogResult{
		Used:      used,
		Threshold: threshold,
		Total:     total,
	}

	if used >= threshold {
		log.Warn().Uint64("used", used).Uint64("threshold", threshold).Msg("Watchdog triggering GC")
		result.ForceGcTook = durationpb.New(forceGC())
		err = fmt.Errorf("usage > threshold, force gc triggered")
	}

	log.Info().Uint64("used", used).Uint64("threshold", threshold).Msg("Memory utilization in bytes")
	return result, err
}
