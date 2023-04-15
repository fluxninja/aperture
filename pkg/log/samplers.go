package log

import (
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// NewRatelimitingSampler return a log sampler that allows max one log entry
// per few seconds.
//
// Use to report cases when we want to warn for visibility, but the warning is
// likely to reoccur, and we do not want to spam logs. Eg. Invalid arguments on
// some intra-aperture RPC call.
//
// This should be preferred over log.Sample(zerolog.Sometimes), so that we can
// be sure that one-off or low-frequency cases are always logged.
// It's preferred to create a new Sampler for each kind of Msg you send.
//
// See also Logger.Autosample()
//
// Note: Return value should be persisted – do not pass directly to
// log.Sample().
func NewRatelimitingSampler() zerolog.Sampler {
	return &zerolog.BurstSampler{
		Burst:  1,
		Period: 5 * time.Minute,
	}
}

// GetAutosampler returns a sampler based on caller's caller location.
//
// Samplers are created using NewRatelimitingSampler().
//
// In most cases, this function should not be used directly, unless implementing
// some other helpers similar to Autosample() or Bug().
func GetAutosampler() zerolog.Sampler {
	// We want separate sampler for every callsite of Autosample() and Bug().
	// Usage of Caller has some runtime cost, but that is acceptable:
	// * for Autosample() a more performant (and verbose) alternative exists,
	// * Bug() in theory never happens, so we should not care too much about perf.
	callerPC, _, _, pcValid := runtime.Caller(2)
	if !pcValid {
		BugWithSampler(global, badCallerSampler).Msg("bug: Cannot get caller info")
		callerPC = 0
	}
	if sampler := getExistingAutosampler(callerPC); sampler != nil {
		return sampler
	}
	autoSamplersLock.Lock()
	defer autoSamplersLock.Unlock()
	// Note: Possible race with another getAutosampler() call, but that is a
	// harmless race.
	sampler := NewRatelimitingSampler()
	autoSamplers[callerPC] = sampler
	return sampler
}

func getExistingAutosampler(callerPC uintptr) zerolog.Sampler {
	autoSamplersLock.RLock()
	defer autoSamplersLock.RUnlock()
	return autoSamplers[callerPC]
}

// BugWithSampler is like Bug() but without auto-sampler logic
//
// In most cases, using Bug() directly should be preferred.
func BugWithSampler(lg *Logger, sampler zerolog.Sampler) *zerolog.Event {
	return lg.Sample(sampler).Warn().Bool(bugKey, true)
}

var (
	autoSamplers     map[uintptr]zerolog.Sampler = make(map[uintptr]zerolog.Sampler)
	autoSamplersLock sync.RWMutex
	badCallerSampler = NewRatelimitingSampler()
)
