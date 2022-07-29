package jobs

import (
	"time"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

// JobStats holds fields to track job statistics.
type JobStats struct {
	Duration time.Duration
}

// JobWatcher is used for tracking completion of Job.
type JobWatcher interface {
	OnJobScheduled()
	OnJobCompleted(status *statusv1.Status, stats JobStats)
}

// JobWatchers is a collection of JobWatcher.
type JobWatchers []JobWatcher

// OnJobScheduled calls OnJobScheduled for each JobWatcher in the collection.
func (jws JobWatchers) OnJobScheduled() {
	for _, jw := range jws {
		jw.OnJobScheduled()
	}
}

// OnJobCompleted calls OnJobCompleted for each JobWatcher in the collection.
func (jws JobWatchers) OnJobCompleted(status *statusv1.Status, jobStats JobStats) {
	for _, jw := range jws {
		jw.OnJobCompleted(status, jobStats)
	}
}

// GroupWatcher is used for tracking completion of JobGroup.
type GroupWatcher interface {
	OnJobRegistered(name string)
	OnJobDeregistered(name string)
	OnJobScheduled(name string)
	OnJobCompleted(name string, status *statusv1.Status, jobStats JobStats)
}

// GroupWatchers is a collection of GroupWatcher.
type GroupWatchers []GroupWatcher

// OnJobRegistered calls OnJobRegistered for each GroupWatcher in the collection.
func (gws GroupWatchers) OnJobRegistered(name string) {
	for _, gw := range gws {
		gw.OnJobRegistered(name)
	}
}

// OnJobDeregistered calls OnJobDeregistered for each GroupWatcher in the collection.
func (gws GroupWatchers) OnJobDeregistered(name string) {
	for _, gw := range gws {
		gw.OnJobDeregistered(name)
	}
}

// OnJobScheduled calls OnJobScheduled for each GroupWatcher in the collection.
func (gws GroupWatchers) OnJobScheduled(name string) {
	for _, gw := range gws {
		gw.OnJobScheduled(name)
	}
}

// OnJobCompleted calls OnJobCompleted for each GroupWatcher in the collection.
func (gws GroupWatchers) OnJobCompleted(name string, status *statusv1.Status, jobStats JobStats) {
	for _, gw := range gws {
		gw.OnJobCompleted(name, status, jobStats)
	}
}
