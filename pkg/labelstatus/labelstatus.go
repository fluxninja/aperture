package labelstatus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// Module is an fx module for providing LabelStatusFactory.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewLabelStatusFactory),
	)
}

// LabelStatusFactory is a factory for creating LabelStatus.
type LabelStatusFactory struct {
	registry status.Registry
}

// NewLabelStatusFactory creates a new LabelStatusFactory.
func NewLabelStatusFactory(statusRegistry status.Registry) *LabelStatusFactory {
	return &LabelStatusFactory{
		registry: statusRegistry.Child("label", "status"),
	}
}

// New creates a new LabelStatus.
func (lsf *LabelStatusFactory) New(labelKey string, policyName string, componentID string) *LabelStatus {
	reg := lsf.registry.Child("label", labelKey)
	return &LabelStatus{
		registry:    reg,
		labelKey:    labelKey,
		policyName:  policyName,
		componentID: componentID,
		timestamp:   time.Time{},
	}
}

// LabelStatus holds the status of the labels.
type LabelStatus struct {
	lock        sync.RWMutex
	registry    status.Registry
	timestamp   time.Time
	labelKey    string
	policyName  string
	componentID string
}

// Setup sets up the LabelsStatus's lifecycle hooks.
func (ls *LabelStatus) Setup(jobGroup *jobs.JobGroup, lifecycle fx.Lifecycle) {
	jobName := fmt.Sprintf("label-status-%s-%s-%s", ls.policyName, ls.componentID, ls.labelKey)
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			job := jobs.NewBasicJob(jobName, ls.setLookupStatus)
			err := jobGroup.RegisterJob(job, jobs.JobConfig{
				ExecutionPeriod: config.MakeDuration(10 * time.Second),
			})
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := jobGroup.DeregisterJob(jobName)
			if err != nil {
				return err
			}
			return nil
		},
	})
}

// SetMissing sets the status to missing with current timestamp.
func (ls *LabelStatus) SetMissing() {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.timestamp = time.Now()
}

func (ls *LabelStatus) setLookupStatus(ctx context.Context) (proto.Message, error) {
	ls.lock.RLock()
	defer ls.lock.RUnlock()

	if ls.timestamp.IsZero() {
		return nil, nil
	}

	labels := map[string]string{
		"policy_name":  ls.policyName,
		"component_id": ls.componentID,
	}

	if time.Since(ls.timestamp) >= 5*time.Minute {
		ls.registry.SetStatus(nil, labels)
		return nil, nil
	} else {
		labels["severity"] = alerts.SeverityCrit.String()
		s := status.NewStatus(nil, errors.New("label "+ls.labelKey+" missing"))
		ls.registry.SetStatus(s, labels)
	}

	return nil, nil
}
