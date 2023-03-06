package podscaler

import (
	"context"
	"encoding/json"
	"path"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sourcegraph/conc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/config"
	discoverykubernetes "github.com/fluxninja/aperture/pkg/discovery/kubernetes"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	etcdwatcher "github.com/fluxninja/aperture/pkg/etcd/watcher"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/paths"
	"github.com/fluxninja/aperture/pkg/status"
)

const podScalerStatusRoot = "pod_scalers"

var fxTag = config.NameTag(podScalerStatusRoot)

// Module returns the fx module for the pod scaler.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideConfigWatcher,
				fx.ResultTags(fxTag),
			),
		),
		fx.Invoke(
			fx.Annotate(
				setupPodScalerFactory,
				fx.ParamTags(
					fxTag,
					discoverykubernetes.FxTag,
					election.FxTag,
				),
			),
		),
	)
}

func provideConfigWatcher(
	etcdClient *etcdclient.Client,
	ai *agentinfo.AgentInfo,
) (notifiers.Watcher, error) {
	agentGroup := ai.GetAgentGroup()

	etcdPath := path.Join(paths.PodScalerConfigPath,
		paths.AgentGroupPrefix(agentGroup))
	watcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

type podScalerFactory struct {
	registry             status.Registry
	decisionsWatcher     notifiers.Watcher
	dynamicConfigWatcher notifiers.Watcher
	controlPointTrackers notifiers.Trackers
	electionTrackers     notifiers.Trackers
	k8sClient            k8s.K8sClient
	etcdClient           *etcdclient.Client
	agentGroup           string
}

// main fx app.
func setupPodScalerFactory(
	watcher notifiers.Watcher,
	controlPointTrackers notifiers.Trackers,
	electionTrackers notifiers.Trackers,
	lifecycle fx.Lifecycle,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
	ai *agentinfo.AgentInfo,
) error {
	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(paths.PodScalerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	dynamicConfigWatcher, err := etcdwatcher.NewWatcher(etcdClient,
		paths.PodScalerDynamicConfigPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", podScalerStatusRoot)
	// logger := reg.GetLogger()

	paFactory := &podScalerFactory{
		controlPointTrackers: controlPointTrackers,
		decisionsWatcher:     decisionsWatcher,
		dynamicConfigWatcher: dynamicConfigWatcher,
		agentGroup:           agentGroup,
		registry:             reg,
		etcdClient:           etcdClient,
		k8sClient:            k8sClient,
		electionTrackers:     electionTrackers,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{paFactory.newPodScalerOptions})
	if err != nil {
		return err
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
func (paFactory *podScalerFactory) newPodScalerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := paFactory.registry.GetLogger()
	wrapperMessage := &policysyncv1.PodScalerWrapper{}
	err := unmarshaller.Unmarshal(wrapperMessage)
	if err != nil || wrapperMessage.PodScaler == nil {
		reg.SetStatus(status.NewStatus(nil, err))
		logger.Warn().Err(err).Msg("Failed to unmarshal pod scaler")
		return fx.Options(), err
	}

	podScalerProto := wrapperMessage.PodScaler
	podScaler := &podScaler{
		Component:        wrapperMessage.GetCommonAttributes(),
		podScalerProto:   podScalerProto,
		registry:         reg,
		podScalerFactory: paFactory,
		scaleWaitGroup:   conc.NewWaitGroup(),
	}
	componentKey := paths.AgentComponentKey(paFactory.agentGroup, podScaler.GetPolicyName(), podScaler.GetComponentId())
	statusEtcdPath := path.Join(paths.PodScalerStatusPath, componentKey)
	podScaler.statusEtcdPath = statusEtcdPath

	return fx.Options(
		fx.Invoke(
			podScaler.setup,
		),
		fx.Supply(
			paFactory.etcdClient,
			fx.Annotate(paFactory.k8sClient, fx.As(new(k8s.K8sClient))),
		),
	), nil
}

// podScaler implement  pod scaler on the agent side.
type podScaler struct {
	stateMutex sync.Mutex
	ctx        context.Context
	k8sClient  k8s.K8sClient
	registry   status.Registry
	iface.Component
	statusWriter      *etcdwriter.Writer
	etcdClient        *etcdclient.Client
	cancel            context.CancelFunc
	scaleCancel       context.CancelFunc
	podScalerFactory  *podScalerFactory
	podScalerProto    *policylangv1.PodScaler
	lastScaleDecision *policysyncv1.ScaleDecision
	scaleWaitGroup    *conc.WaitGroup
	controlPoint      discoverykubernetes.AutoscaleControlPoint
	statusEtcdPath    string
	dryRun            bool
	isLeader          bool
}

func (pa *podScaler) setup(
	lifecycle fx.Lifecycle,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
) error {
	logger := pa.registry.GetLogger()
	pa.etcdClient = etcdClient
	pa.k8sClient = k8sClient
	etcdKey := paths.AgentComponentKey(pa.podScalerFactory.agentGroup,
		pa.GetPolicyName(),
		pa.GetComponentId())

	// election notifier
	electionNotifier := notifiers.NewBasicKeyNotifier(election.ElectionResultKey, pa.electionResultCallback)

	// decision notifier
	decisionUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaler,
		pa.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}
	// dynamic config notifier
	dynamicConfigUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	dynamicConfigNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		dynamicConfigUnmarshaler,
		pa.dynamicConfigUpdateCallback,
	)
	if err != nil {
		return err
	}
	// control point notifier
	// read the configured control point from the horizontal pod scaler proto
	controlPointSelector := pa.podScalerProto.KubernetesObjectSelector
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
	controlPointNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(key),
		controlPointUnmarshaler,
		pa.controlPointUpdateCallback,
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var err error
			pa.statusWriter = etcdwriter.NewWriter(pa.etcdClient, true)
			pa.ctx, pa.cancel = context.WithCancel(context.Background())
			scaleActuatorProto := pa.podScalerProto.GetScaleActuator()
			if scaleActuatorProto != nil {
				pa.updateDynamicConfig(scaleActuatorProto.GetDefaultConfig())
			}
			// add election notifier
			err = pa.podScalerFactory.electionTrackers.AddKeyNotifier(electionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add election notifier")
				return err
			}
			// add decisions notifier
			err = pa.podScalerFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}
			// add dynamic config notifier
			err = pa.podScalerFactory.dynamicConfigWatcher.AddKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add dynamic config notifier")
			}
			// add control point notifier
			err = pa.podScalerFactory.controlPointTrackers.AddKeyNotifier(controlPointNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add control point notifier")
			}

			return err
		},
		OnStop: func(ctx context.Context) error {
			var merr, err error
			// remove election notifier
			err = pa.podScalerFactory.electionTrackers.RemoveKeyNotifier(electionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove election notifier")
				merr = multierror.Append(merr, err)
			}
			// remove dynamic config notifier
			err = pa.podScalerFactory.dynamicConfigWatcher.RemoveKeyNotifier(dynamicConfigNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove dynamic config notifier")
				merr = multierror.Append(merr, err)
			}
			// remove decisions notifier
			err = pa.podScalerFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierror.Append(merr, err)
			}
			// remove control point notifier
			err = pa.podScalerFactory.controlPointTrackers.RemoveKeyNotifier(controlPointNotifier)
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

func (pa *podScaler) electionResultCallback(_ notifiers.Event) {
	log.Info().Msg("Election result callback")

	// invoke the lastScaleDecision
	pa.stateMutex.Lock()
	defer pa.stateMutex.Unlock()
	if pa.lastScaleDecision != nil {
		pa.scale(pa.lastScaleDecision)
	}
	pa.isLeader = true
}

func (pa *podScaler) dynamicConfigUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := pa.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Dynamic config removed")
		// revert to default config
		scaleActuatorProto := pa.podScalerProto.GetScaleActuator()
		if scaleActuatorProto != nil {
			pa.updateDynamicConfig(scaleActuatorProto.GetDefaultConfig())
		}
		return
	}

	var wrapperMessage policysyncv1.PodScalerDynamicConfigWrapper
	err := unmarshaller.Unmarshal(&wrapperMessage)
	if err != nil || wrapperMessage.PodScalerDynamicConfig == nil {
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
	dynamicConfig := wrapperMessage.PodScalerDynamicConfig
	pa.updateDynamicConfig(dynamicConfig)
}

func (pa *podScaler) updateDynamicConfig(dynamicConfig *policylangv1.PodScaler_ScaleActuator_DynamicConfig) {
	if dynamicConfig == nil {
		pa.dryRun = false
		return
	}
	pa.dryRun = dynamicConfig.GetDryRun()
}

func (pa *podScaler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	pa.stateMutex.Lock()
	defer pa.stateMutex.Unlock()
	logger := pa.registry.GetLogger()
	pa.lastScaleDecision = nil

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
	if !pa.dryRun {
		pa.lastScaleDecision = scaleDecision
		if pa.isLeader {
			pa.scale(scaleDecision)
		}
	}
}

// scale scales the associated Kubernetes object. NOTE: not thread safe, needs to be called under podScaler.scaleMutex.
func (pa *podScaler) scale(scaleDecision *policysyncv1.ScaleDecision) {
	// Take mutex to prevent concurrent scale operations
	replicas := scaleDecision.GetDesiredReplicas()

	// Cancel any existing scale operation
	if pa.scaleCancel != nil {
		pa.scaleCancel()
	}
	// Cancel any existing scale operation
	ctx, cancel := context.WithCancel(pa.ctx)
	pa.scaleCancel = cancel
	// Wait on existing scaleWaitGroup to make sure previous scale operation is complete
	pa.scaleWaitGroup.Wait()
	// Create a new scaleWaitGroup
	pa.scaleWaitGroup = conc.NewWaitGroup()
	pa.scaleWaitGroup.Go(func() {
		cp := pa.controlPoint
		targetGK := schema.GroupKind{
			Group: cp.Group,
			Kind:  cp.Kind,
		}

		operation := func() error {
			scale, targetGR, err := pa.k8sClient.ScaleForGroupKind(ctx, cp.Namespace, cp.Name, targetGK)
			if err != nil {
				// TODO: update status
				log.Error().Err(err).Msgf("Unable to get scale for %v", cp)
				return err
			}

			if scale.Spec.Replicas != replicas {
				scale.Spec.Replicas = replicas
				_, err = pa.k8sClient.GetScaleClient().Scales(cp.Namespace).Update(ctx, targetGR, scale, metav1.UpdateOptions{})
				if err != nil {
					// TODO: update status
					log.Error().Err(err).Msg("Unable to update scale subresource")
					return err
				}
			}
			return nil
		}

		merr := backoff.Retry(operation, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
		if merr != nil {
			log.Error().Err(merr).Msgf("Context canceled while invoking scale for %v", cp)
		}
	})
}

func (pa *podScaler) controlPointUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
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
			PolicyName:  pa.GetPolicyName(),
			PolicyHash:  pa.GetPolicyHash(),
			ComponentId: pa.GetComponentId(),
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
