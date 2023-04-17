// +kubebuilder:validation:Optional
package watchdogconfig

import (
	"sort"

	jobs "github.com/fluxninja/aperture/pkg/jobs/config"
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
	// Minimum GoGC sets the minimum garbage collection target percentage for heap driven Watchdogs. This setting helps avoid over scheduling.
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

// WatermarksPolicy creates a Watchdog policy that schedules GC at concrete watermarks.
// swagger:model
// +kubebuilder:object:generate=true
type WatermarksPolicy struct {
	// Watermarks are increasing limits on which to trigger GC. Watchdog disarms when the last watermark is surpassed. It's recommended to set an extreme watermark for the last element (for example, 0.99).
	Watermarks []float64 `json:"watermarks,omitempty" validate:"omitempty,dive,gte=0,lte=1" default:"[0.50,0.75,0.80,0.85,0.90,0.95,0.99]"`

	// internal fields
	// swagger:ignore
	thresholds []uint64

	PolicyCommon `json:",inline"`
}

// NextThreshold implements Policy.
func (policy *WatermarksPolicy) NextThreshold(total, used uint64) uint64 {
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
	return PolicyTempDisabled
}

// AdaptivePolicy creates a policy that forces GC when the usage surpasses the configured factor of the available memory. This policy calculates next target as usage+(limit-usage)*factor.
// swagger:model
// +kubebuilder:object:generate=true
type AdaptivePolicy struct {
	PolicyCommon `json:",inline"`

	// Factor sets user-configured limit of available memory
	Factor float64 `json:"factor" validate:"gte=0,lte=1" default:"0.50"`
}

// NextThreshold implements Policy.
func (policy *AdaptivePolicy) NextThreshold(total, used uint64) uint64 {
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
