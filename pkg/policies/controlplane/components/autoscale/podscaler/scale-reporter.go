package podscaler

import (
	"context"
	"sync"

	policyprivatev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/private/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/config"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
	"go.uber.org/fx"
)

// fxTag is PodScaler's Status Watcher's Fx Tag.
var fxTag = config.NameTag("scale_status_watcher")

// scaleReporterModule returns the fx options for pod scaler in the main app.
func scaleReporterModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideWatcher,
				fx.ResultTags(fxTag),
			),
			fx.Annotate(
				provideFxOptionsFunc,
				fx.ParamTags(fxTag),
				fx.ResultTags(iface.FxOptionsFuncTag),
			),
		),
	)
}

func provideWatcher(
	etcdClient *etcdclient.Client,
	lc fx.Lifecycle,
) (notifiers.Watcher, error) {
	etcdPath := paths.PodScalerStatusPath
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}
	notifiers.WatcherLifecycle(lc, watcher, nil)

	return watcher, nil
}

func provideFxOptionsFunc(watcher notifiers.Watcher) notifiers.FxOptionsFunc {
	return func(key notifiers.Key, _ config.Unmarshaller, _ status.Registry) (fx.Option, error) {
		return fx.Supply(fx.Annotate(
			watcher, fx.ResultTags(fxTag),
		)), nil
	}
}

// ScaleReporter struct.
type ScaleReporter struct {
	lock          sync.RWMutex
	policyReadAPI iface.Policy
	// statusRegistry status.Registry
	scaleStatus *policysyncv1.ScaleStatus
	etcdKey     string
	agentGroup  string
}

// Make sure ScaleReporter implements runtime.Component.
var _ runtime.Component = (*ScaleReporter)(nil)

// Name implements runtime.Component.Name.
func (*ScaleReporter) Name() string { return "ScaleReporter" }

// Type implements runtime.Component.Type.
func (sr *ScaleReporter) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.ShortDescription.
func (sr *ScaleReporter) ShortDescription() string { return sr.agentGroup }

// IsActuator implements runtime.Component.
func (*ScaleReporter) IsActuator() bool { return false }

// NewScaleReporterAndOptions returns a new ScaleReporter and its fx options.
func NewScaleReporterAndOptions(
	podScaleReporterProto *policyprivatev1.PodScaleReporter,
	_ runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, fx.Option, error) {
	agentGroup := podScaleReporterProto.GetAgentGroup()
	podScalerComponentID := podScaleReporterProto.GetPodScalerComponentId()
	etcdKey := paths.AgentComponentKey(agentGroup, policyReadAPI.GetPolicyName(), podScalerComponentID)
	sr := &ScaleReporter{
		policyReadAPI: policyReadAPI,
		etcdKey:       etcdKey,
		agentGroup:    agentGroup,
	}

	return sr, fx.Options(
		fx.Invoke(
			fx.Annotate(
				sr.setupWatch,
				fx.ParamTags(fxTag),
			),
		),
	), nil
}

func (sr *ScaleReporter) setupWatch(
	watcher notifiers.Watcher,
	lifecycle fx.Lifecycle,
) error {
	statusUnmarshaller, protoErr := config.NewProtobufUnmarshaller(nil)
	if protoErr != nil {
		return protoErr
	}

	// status notifier
	statusNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(sr.etcdKey),
		statusUnmarshaller,
		sr.statusUpdateCallback,
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := watcher.AddKeyNotifier(statusNotifier)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := watcher.RemoveKeyNotifier(statusNotifier)
			if err != nil {
				return err
			}
			return nil
		},
	})
	return nil
}

func (sr *ScaleReporter) statusUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	sr.lock.Lock()
	defer sr.lock.Unlock()
	logger := sr.policyReadAPI.GetStatusRegistry().GetLogger()
	if event.Type == notifiers.Remove {
		logger.Info().Msg("ScaleReporter: status removed")
		sr.scaleStatus = nil
		return
	}

	var wrapperMessage policysyncv1.ScaleStatusWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	scaleStatus := wrapperMessage.GetScaleStatus()
	if err != nil || scaleStatus == nil {
		logger.Warn().Err(err).Msg("Failed to unmarshal config wrapper")
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Failed to get common attributes from ScaleStatusWrapper")
		return
	}
	if commonAttributes.PolicyHash != sr.policyReadAPI.GetPolicyHash() {
		return
	}

	sr.scaleStatus = scaleStatus
}

// Execute implements runtime.Component.
func (sr *ScaleReporter) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	outPortReadings := make(runtime.PortToReading)
	configuredReplicasReading := runtime.InvalidReading()
	actualReplicasReading := runtime.InvalidReading()

	sr.lock.RLock()
	defer sr.lock.RUnlock()
	if sr.scaleStatus != nil {
		configuredReplicasReading = runtime.NewReading(float64(sr.scaleStatus.ConfiguredReplicas))
		actualReplicasReading = runtime.NewReading(float64(sr.scaleStatus.ActualReplicas))
	}

	outPortReadings["configured_replicas"] = []runtime.Reading{configuredReplicasReading}
	outPortReadings["actual_replicas"] = []runtime.Reading{actualReplicasReading}

	return outPortReadings, nil
}

// DynamicConfigUpdate is a no-op for this component.
func (sr *ScaleReporter) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
}
