package podautoscaler

import (
	"context"
	"encoding/json"
	"path"
	"sync"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	discoverykubernetes "github.com/fluxninja/aperture/pkg/discovery/kubernetes"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const podAutoscalerStatusRoot = "pod_autoscalers"

var fxNameTag = config.NameTag(podAutoscalerStatusRoot)

// Module returns the fx module for the pod autoscaler.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideWatcher,
				fx.ResultTags(fxNameTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupPodAutoscalerFactory,
				fx.ParamTags(
					fxNameTag,
					config.NameTag("kubernetes_control_points"),
				),
			),
		),
	)
}

func provideWatcher(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroup := ai.GetAgentGroup()

	etcdPath := path.Join(paths.PodAutoscalerConfigPath,
		paths.AgentGroupPrefix(agentGroup))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type podAutoscalerFactory struct {
	registry             status.Registry
	decisionsWatcher     notifiers.Watcher
	dynamicConfigWatcher notifiers.Watcher
	controlPointTrackers notifiers.Trackers
	etcdClient           *etcdclient.Client
	k8sClient            k8s.K8sClient
	agentGroup           string
}

// main fx app.
func setupPodAutoscalerFactory(
	watcher notifiers.Watcher,
	controlPointTrackers notifiers.Trackers,
	lifecycle fx.Lifecycle,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
	ai *agentinfo.AgentInfo,
) error {
	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(paths.PodAutoscalerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		paths.PodAutoscalerDynamicConfigPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child(podAutoscalerStatusRoot)
	// logger := reg.GetLogger()

	paFactory := &podAutoscalerFactory{
		controlPointTrackers: controlPointTrackers,
		decisionsWatcher:     decisionsWatcher,
		dynamicConfigWatcher: dynamicConfigWatcher,
		agentGroup:           agentGroup,
		registry:             reg,
		etcdClient:           etcdClient,
		k8sClient:            k8sClient,
	}

	fxDriver := &notifiers.FxDriver{
		FxOptionsFuncs: []notifiers.FxOptionsFunc{
			paFactory.newPodAutoscalerOptions,
		},
		UnmarshalPrefixNotifier: notifiers.UnmarshalPrefixNotifier{
			GetUnmarshallerFunc: config.NewProtobufUnmarshaller,
		},
		StatusRegistry:     reg,
		PrometheusRegistry: prometheusRegistry,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := decisionsWatcher.Start()
			if err != nil {
				return err
			}
			err = dynamicConfigWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var err, merr error
			err = dynamicConfigWatcher.Stop()
			if err != nil {
				merr = multierror.Append(merr, err)
			}
			err = decisionsWatcher.Stop()
			if err != nil {
				merr = multierror.Append(merr, err)
			}
			return merr
		},
	})

	notifiers.WatcherLifecycle(lifecycle, watcher, []notifiers.PrefixNotifier{fxDriver})

	return nil
}

// per component fx app.
func (paFactory *podAutoscalerFactory) newPodAutoscalerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := paFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.PodAutoscalerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.PodAutoscaler == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal pod autoscaler")
		return fx.Options(), err
	}

	podAutoscalerProto := wrapperMessage.PodAutoscaler
	podAutoscaler := &podAutoscaler{
		Component:            wrapperMessage.GetCommonAttributes(),
		podAutoscalerProto:   podAutoscalerProto,
		registry:             reg,
		podAutoscalerFactory: paFactory,
	}
	componentKey := paths.AgentComponentKey(paFactory.agentGroup, podAutoscaler.GetPolicyName(), int64(podAutoscaler.GetComponentIndex()))
	statusEtcdPath := path.Join(paths.PodAutoscalerStatusPath, componentKey)
	podAutoscaler.statusEtcdPath = statusEtcdPath

	return fx.Options(
		fx.Invoke(
			fx.Annotate(
				podAutoscaler.setup,
				fx.ParamTags("kubernetes_control_points"),
			),
		),
		fx.Supply(
			paFactory.etcdClient,
			paFactory.k8sClient,
		),
	), nil
}

// podAutoscaler implement pod auto scaler on the agent side.
type podAutoscaler struct {
	scaleMutex sync.Mutex
	ctx        context.Context
	k8sClient  k8s.K8sClient
	registry   status.Registry
	iface.Component
	podAutoscalerProto   *policylangv1.PodAutoscaler
	podAutoscalerFactory *podAutoscalerFactory
	etcdClient           *etcdclient.Client
	scaleCancel          context.CancelFunc
	cancel               context.CancelFunc
	statusWriter         *etcdwriter.Writer
	controlPoint         discoverykubernetes.ControlPoint
	statusEtcdPath       string
	dryRun               bool
}

func (pa *podAutoscaler) setup(
	lifecycle fx.Lifecycle,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
) error {
	logger := pa.registry.GetLogger()
	pa.etcdClient = etcdClient
	pa.k8sClient = k8sClient
	etcdKey := paths.AgentComponentKey(pa.podAutoscalerFactory.agentGroup,
		pa.GetPolicyName(),
		pa.GetComponentIndex())
	// decision notifier
	decisionUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaler,
		pa.decisionUpdateCallback,
	)
	// dynamic config notifier
	dynamicConfigUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	dynamicConfigNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		dynamicConfigUnmarshaler,
		pa.dynamicConfigUpdateCallback,
	)
	// control point notifier
	// read the configured control point from the pod autoscaler proto
	controlPointSelector := pa.podAutoscalerProto.KubernetesObjectSelector
	controlPoint, err := discoverykubernetes.ControlPointFromSelector(controlPointSelector)
	if err != nil {
		return err
	}
	pa.controlPoint = controlPoint
	key, keyErr := json.Marshal(controlPoint)
	if keyErr != nil {
		return keyErr
	}
	// control point notifier
	controlPointUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	controlPointNotifier := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(key),
		controlPointUnmarshaler,
		pa.controlPointUpdateCallback,
	)

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var err error
			pa.statusWriter = etcdwriter.NewWriter(pa.etcdClient, true)
			pa.ctx, pa.cancel = context.WithCancel(context.Background())
			scaleActuatorProto := pa.podAutoscalerProto.GetScaleActuator()
			if scaleActuatorProto != nil {
				pa.updateDynamicConfig(scaleActuatorProto.GetDefaultConfig())
			}
			// add decisions notifier
			err = pa.podAutoscalerFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}
			// add dynamic config notifier
			err = pa.podAutoscalerFactory.dynamicConfigWatcher.AddKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add dynamic config notifier")
			}
			// add control point notifier
			err = pa.podAutoscalerFactory.controlPointTrackers.AddKeyNotifier(controlPointNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add control point notifier")
			}

			return err
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			// remove dynamic config notifier
			err = pa.podAutoscalerFactory.dynamicConfigWatcher.RemoveKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove dynamic config notifier")
				merr = multierror.Append(merr, err)
			}
			// remove decisions notifier
			err = pa.podAutoscalerFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierror.Append(merr, err)
			}
			// remove control point notifier
			err = pa.podAutoscalerFactory.controlPointTrackers.RemoveKeyNotifier(controlPointNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove control point notifier")
				merr = multierror.Append(merr, err)
			}
			pa.registry.SetStatus(status.NewStatus(nil, merr))
			pa.cancel()
			pa.statusWriter.Close()
			_, err = pa.etcdClient.KV.Delete(clientv3.WithRequireLeader(ctx), pa.statusEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scale status")
				merr = multierr.Append(merr, err)
			}
			pa.registry.SetStatus(status.NewStatus(nil, merr))
			return merr
		},
	})
	return nil
}

func (pa *podAutoscaler) dynamicConfigUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := pa.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Dynamic config removed")
		// revert to default config
		scaleActuatorProto := pa.podAutoscalerProto.GetScaleActuator()
		if scaleActuatorProto != nil {
			pa.updateDynamicConfig(scaleActuatorProto.GetDefaultConfig())
		}
		return
	}

	var wrapperMessage policysyncv1.PodAutoscalerDynamicConfigWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.PodAutoscalerDynamicConfig == nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Common attributes not found")
		return
	}
	if commonAttributes.PolicyHash != pa.GetPolicyHash() {
		return
	}
	dynamicConfig := wrapperMessage.PodAutoscalerDynamicConfig
	pa.updateDynamicConfig(dynamicConfig)
}

func (pa *podAutoscaler) updateDynamicConfig(dynamicConfig *policylangv1.PodAutoscaler_ScaleActuator_DynamicConfig) {
	if dynamicConfig == nil {
		pa.dryRun = false
		return
	}
	pa.dryRun = dynamicConfig.GetDryRun()
}

func (pa *podAutoscaler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := pa.registry.GetLogger()
	if pa.dryRun {
		return
	}
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Decision removed")
		return
	}

	var wrapperMessage policysyncv1.ScaleDecisionWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil {
		return
	}
	commonAttributes := wrapperMessage.GetCommonAttributes()
	if commonAttributes == nil {
		log.Error().Msg("Decision missing common attributes")
		return
	}
	if commonAttributes.PolicyHash != pa.GetPolicyHash() {
		return
	}
	scaleDecision := wrapperMessage.ScaleDecision
	pa.scale(scaleDecision.DesiredReplicas)
}

// scale scales the associated Kubernetes object.
func (pa *podAutoscaler) scale(replicas int32) {
	// Take mutex to prevent concurrent scale operations
	pa.scaleMutex.Lock()
	defer pa.scaleMutex.Unlock()
	// Cancel any existing scale operation
	if pa.scaleCancel != nil {
		pa.scaleCancel()
	}
	ctx, cancel := context.WithCancel(pa.ctx)
	pa.scaleCancel = cancel
	panichandler.Go(func() {
		cp := pa.controlPoint
		targetGK := schema.GroupKind{
			Group: cp.Group,
			Kind:  cp.Kind,
		}

		scale, targetGR, err := pa.k8sClient.ScaleForGroupKind(ctx, cp.Namespace, cp.Name, targetGK)
		if err != nil {
			// TODO: update status
			log.Error().Err(err).Msgf("Unable to get scale for %v", cp)
			return
		}

		if scale.Spec.Replicas != replicas {
			_, err = pa.k8sClient.GetScaleClient().Scales(cp.Namespace).Update(ctx, targetGR, scale, metav1.UpdateOptions{})
			if err != nil {
				// TODO: update status
				log.Error().Err(err).Msg("Unable to update scale subresource")
				return
			}
		}
	})
}

func (pa *podAutoscaler) controlPointUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := pa.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Control point removed")
		pa.statusWriter.Delete(pa.statusEtcdPath)
		return
	}

	var scaleStatus policysyncv1.ScaleStatus
	err := unmarshaller.Unmarshal(&scaleStatus)
	if err != nil {
		// TODO: update status
		log.Error().Err(err).Msg("Unable to unmarshal scale status")
		return
	}

	// create a wrapper message
	wrapperMessage := policysyncv1.ScaleStatusWrapper{
		CommonAttributes: &policysyncv1.CommonAttributes{
			PolicyName:     pa.GetPolicyName(),
			PolicyHash:     pa.GetPolicyHash(),
			ComponentIndex: pa.GetComponentIndex(),
		},
		ScaleStatus: &scaleStatus,
	}

	// marshal the wrapper message
	data, err := proto.Marshal(&wrapperMessage)
	if err != nil {
		// TODO: update status
		log.Error().Err(err).Msg("Unable to marshal scale status")
		return
	}
	pa.statusWriter.Write(pa.statusEtcdPath, data)
}
