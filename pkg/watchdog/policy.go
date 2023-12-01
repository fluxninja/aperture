package watchdog

import (
	"context"
	"runtime"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	watchdogv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/watchdog/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	watchdogconfig "github.com/fluxninja/aperture/v2/pkg/watchdog/config"
)

type usageFn func() (total uint64, usage uint64, err error)

func forceGC() time.Duration {
	log.Info().Msg("watchdog is forcing GC")
	start := time.Now()
	runtime.GC()
	took := time.Since(start)
	log.Info().Dur("took", took).Msg("Watchdog triggered GC finished")
	return took
}

func check(policy watchdogconfig.Policy, _ context.Context, fn usageFn) (proto.Message, error) {
	total, used, err := fn()
	if err != nil {
		log.Warn().Err(err).Msg("failed to get memory usage")
		return nil, nil
	}

	threshold := policy.NextThreshold(total, used)

	result := &watchdogv1.WatchdogResult{
		Used:      used,
		Threshold: threshold,
		Total:     total,
	}

	if used >= threshold {
		log.Warn().Uint64("used", used).Uint64("threshold", threshold).Msg("Watchdog triggering GC")
		result.ForceGcTook = durationpb.New(forceGC())
		log.Warn().Msg("usage > threshold, force gc triggered")
	}

	log.Info().Uint64("used", used).Uint64("threshold", threshold).Msg("Memory utilization in bytes")
	return result, nil
}
