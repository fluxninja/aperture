package prometheus

import (
	"context"
	"time"

	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusmodel "github.com/prometheus/common/model"
	"google.golang.org/protobuf/proto"

	"github.com/FluxNinja/aperture/pkg/jobs"
)

// AlertCallback is a callback function that gets invoked respectively when the alert gets active, inactive.
type AlertCallback func(context.Context, ...interface{}) (proto.Message, error)

type alertState int64

const (
	inactive alertState = iota
	pending
	active
)

type alertQuery struct {
	activeAt              time.Time
	savedError            error
	savedDetails          proto.Message
	alertActiveCallback   AlertCallback
	alertInactiveCallback AlertCallback
	state                 alertState
	forDuration           time.Duration
}

// NewAlertQueryJob takes Alert active and Alert inactive callbacks which get invoked when the alert gets active, inactive respectively.
// Also, it takes an error callback which gets invoked when there's an error from running PromQL.
// Alert is computed via a PromQL query using semantics similar to Prometheus alert rules.
// It returns a callback compatible with scheduler BasicJob.
func NewAlertQueryJob(
	query string,
	endTimestamp time.Time,
	promAPI prometheusv1.API,
	timeout time.Duration,
	forDuration time.Duration,
	alertActiveCallback,
	alertInactiveCallback AlertCallback,
	errorCallback PromErrorCallback,
	cbArgs ...interface{},
) jobs.JobCallback {
	aq := &alertQuery{forDuration: forDuration, alertActiveCallback: alertActiveCallback, alertInactiveCallback: alertInactiveCallback}
	return NewPromQueryJob(query, endTimestamp, promAPI, timeout, aq.execute, errorCallback, cbArgs...)
}

func (aq *alertQuery) execute(jobCtxt context.Context, value prometheusmodel.Value, cbArgs ...interface{}) (proto.Message, error) {
	activeNow := false
	if _, ok := value.(*prometheusmodel.Scalar); ok {
		activeNow = true
	} else if vector, ok := value.(prometheusmodel.Vector); ok {
		if vector.Len() > 0 {
			activeNow = true
		}
	} else if matrix, ok := value.(prometheusmodel.Matrix); ok {
		if matrix.Len() > 0 {
			activeNow = true
		}
	} else if _, ok := value.(*prometheusmodel.String); ok {
		activeNow = true
	}

	if aq.state == inactive {
		if activeNow {
			aq.activeAt = time.Now()
			if aq.forDuration == 0 {
				// Transition
				aq.state = active
				aq.savedDetails, aq.savedError = aq.alertActiveCallback(jobCtxt, cbArgs...)
			}
		}
	} else if aq.state == pending {
		if activeNow {
			// Make sure it is active for forDuration before marking as active
			if time.Since(aq.activeAt) >= aq.forDuration {
				// Transition
				aq.state = active
				aq.savedDetails, aq.savedError = aq.alertActiveCallback(jobCtxt, cbArgs...)
			}
		} else {
			aq.state = inactive
		}
	} else if aq.state == active {
		if !activeNow {
			// Transition
			aq.state = inactive
			aq.savedDetails, aq.savedError = aq.alertInactiveCallback(jobCtxt, cbArgs...)
		}
	}

	return aq.savedDetails, aq.savedError
}
