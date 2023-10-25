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
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	etcdwatcher "github.com/fluxninja/aperture/v2/pkg/etcd/watcher"
	"github.com/fluxninja/aperture/v2/pkg/k8s"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	autoscalek8sconfig "github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/config"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/status"
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
					discovery.FxTag,
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
	controlPointTrackers notifiers.Trackers
	k8sClient            k8s.K8sClient
	etcdClient           *etcdclient.Client
	agentGroup           string
}

// main fx app.
func setupPodScalerFactory(
	watcher notifiers.Watcher,
	controlPointTrackers notifiers.Trackers,
	lifecycle fx.Lifecycle,
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
	ai *agentinfo.AgentInfo,
	cfg autoscalek8sconfig.AutoScaleKubernetesConfig,
) error {
	if !cfg.Enabled {
		log.Info().Msg("Kubernetes AutoScaler is disabled")
		return nil
	}
	if k8sClient == nil {
		log.Info().Msg("Not in Kubernetes cluster, omitting AutoScaler")
		return nil
	}

	agentGroup := ai.GetAgentGroup()
	etcdPath := path.Join(paths.PodScalerDecisionsPath)
	decisionsWatcher, err := etcdwatcher.NewWatcher(etcdClient, etcdPath)
	if err != nil {
		return err
	}

	reg := statusRegistry.Child("component", podScalerStatusRoot)
	// logger := reg.GetLogger()

	psFactory := &podScalerFactory{
		controlPointTrackers: controlPointTrackers,
		decisionsWatcher:     decisionsWatcher,
		agentGroup:           agentGroup,
		registry:             reg,
		etcdClient:           etcdClient,
		k8sClient:            k8sClient,
	}

	fxDriver, err := notifiers.NewFxDriver(reg, prometheusRegistry,
		config.NewProtobufUnmarshaller,
		[]notifiers.FxOptionsFunc{psFactory.newPodScalerOptions})
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := decisionsWatcher.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var err, merr error
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
func (psFactory *podScalerFactory) newPodScalerOptions(
	key notifiers.Key,
	unmarshaller config.Unmarshaller,
	reg status.Registry,
) (fx.Option, error) {
	logger := psFactory.registry.GetLogger()
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
		podScalerFactory: psFactory,
		scaleWaitGroup:   conc.NewWaitGroup(),
	}
	componentKey := paths.AgentComponentKey(psFactory.agentGroup, podScaler.GetPolicyName(), podScaler.GetComponentId())
	statusEtcdPath := path.Join(paths.PodScalerStatusPath, componentKey)
	podScaler.statusEtcdPath = statusEtcdPath

	return fx.Options(
		fx.Invoke(
			podScaler.setup,
		),
		fx.Supply(
			psFactory.etcdClient,
			fx.Annotate(psFactory.k8sClient, fx.As(new(k8s.K8sClient))),
		),
	), nil
}

// podScaler implement  pod scaler on the agent side.
type podScaler struct {
	iface.Component
	ctx               context.Context
	k8sClient         k8s.K8sClient
	registry          status.Registry
	cancel            context.CancelFunc
	etcdClient        *etcdclient.Client
	scaleCancel       context.CancelFunc
	podScalerFactory  *podScalerFactory
	podScalerProto    *policylangv1.PodScaler
	lastScaleDecision *policysyncv1.ScaleDecision
	scaleWaitGroup    *conc.WaitGroup
	controlPoint      discovery.AutoScaleControlPoint
	statusEtcdPath    string
	stateMutex        sync.Mutex
}

// podScaler implements etcdclient.ElectionWatcher.
var _ etcdclient.ElectionWatcher = (*podScaler)(nil)

func (ps *podScaler) setup(
	lifecycle fx.Lifecycle,
	etcdClient *etcdclient.Client,
	k8sClient k8s.K8sClient,
) error {
	logger := ps.registry.GetLogger()
	ps.etcdClient = etcdClient
	ps.k8sClient = k8sClient
	etcdKey := paths.AgentComponentKey(ps.podScalerFactory.agentGroup,
		ps.GetPolicyName(),
		ps.GetComponentId())

	// decision notifier
	decisionUnmarshaler, err := config.NewProtobufUnmarshaller(nil)
	if err != nil {
		return err
	}
	decisionNotifier, err := notifiers.NewUnmarshalKeyNotifier(
		notifiers.Key(etcdKey),
		decisionUnmarshaler,
		ps.decisionUpdateCallback,
	)
	if err != nil {
		return err
	}
	// control point notifier
	// read the configured control point from the horizontal pod scaler proto
	controlPointSelector := ps.podScalerProto.KubernetesObjectSelector
	controlPoint, err := discovery.ControlPointFromSelector(controlPointSelector)
	if err != nil {
		return err
	}
	ps.controlPoint = controlPoint
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
		ps.controlPointUpdateCallback,
	)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var err error
			ps.ctx, ps.cancel = context.WithCancel(context.Background())

			// add decisions notifier
			err = ps.podScalerFactory.decisionsWatcher.AddKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add decision notifier")
				return err
			}
			// add control point notifier
			err = ps.podScalerFactory.controlPointTrackers.AddKeyNotifier(controlPointNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to add control point notifier")
			}

			return err
		},
		OnStop: func(_ context.Context) error {
			var merr, err error

			// remove decisions notifier
			err = ps.podScalerFactory.decisionsWatcher.RemoveKeyNotifier(decisionNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove decision notifier")
				merr = multierror.Append(merr, err)
			}
			// remove control point notifier
			err = ps.podScalerFactory.controlPointTrackers.RemoveKeyNotifier(controlPointNotifier)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to remove control point notifier")
				merr = multierror.Append(merr, err)
			}
			ps.registry.SetStatus(status.NewStatus(nil, merr))
			ps.cancel()
			ps.etcdClient.Delete(ps.statusEtcdPath)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to delete scale status")
				merr = multierr.Append(merr, err)
			}
			ps.registry.SetStatus(status.NewStatus(nil, merr))
			return merr
		},
	})

	ps.etcdClient.AddElectionWatcher(ps)

	return nil
}

// OnLeaderStart is called when the agent becomes the leader.
func (ps *podScaler) OnLeaderStart() {
	ps.stateMutex.Lock()
	defer ps.stateMutex.Unlock()
	if ps.lastScaleDecision != nil {
		ps.scale(ps.lastScaleDecision)
	}
}

// OnLeaderStop is called when the agent stops being the leader.
func (ps *podScaler) OnLeaderStop() {
	ps.stateMutex.Lock()
	defer ps.stateMutex.Unlock()
	if ps.scaleCancel != nil {
		ps.scaleCancel()
		ps.scaleCancel = nil
	}
}

func (ps *podScaler) decisionUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	ps.stateMutex.Lock()
	defer ps.stateMutex.Unlock()
	logger := ps.registry.GetLogger()
	ps.lastScaleDecision = nil

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
	if commonAttributes.PolicyHash != ps.GetPolicyHash() {
		return
	}
	scaleDecision := wrapperMessage.ScaleDecision
	ps.lastScaleDecision = scaleDecision
	if ps.podScalerFactory.etcdClient.IsLeader() {
		ps.scale(scaleDecision)
	}
}

// scale scales the associated Kubernetes object. NOTE: not thread safe, needs to be called under podScaler.scaleMutex.
func (ps *podScaler) scale(scaleDecision *policysyncv1.ScaleDecision) {
	// This is called under stateMutex to prevent concurrent scale operations
	replicas := scaleDecision.GetDesiredReplicas()

	// Cancel any existing scale operation
	if ps.scaleCancel != nil {
		ps.scaleCancel()
	}
	// Cancel any existing scale operation
	ctx, cancel := context.WithCancel(ps.ctx)
	ps.scaleCancel = cancel
	// Wait on existing scaleWaitGroup to make sure previous scale operation is complete
	ps.scaleWaitGroup.Wait()
	// Create a new scaleWaitGroup
	ps.scaleWaitGroup = conc.NewWaitGroup()
	ps.scaleWaitGroup.Go(func() {
		cp := ps.controlPoint
		targetGK := schema.GroupKind{
			Group: cp.Group,
			Kind:  cp.Kind,
		}

		operation := func() error {
			scale, targetGR, err := ps.k8sClient.ScaleForGroupKind(ctx, cp.Namespace, cp.Name, targetGK)
			if err != nil {
				// TODO: update status
				log.Error().Err(err).Msgf("Unable to get scale for %v", cp)
				return err
			}

			if scale.Spec.Replicas != replicas {
				scale.Spec.Replicas = replicas
				_, err = ps.k8sClient.GetScaleClient().Scales(cp.Namespace).Update(ctx, targetGR, scale, metav1.UpdateOptions{})
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

func (ps *podScaler) controlPointUpdateCallback(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := ps.registry.GetLogger()
	if event.Type == notifiers.Remove {
		logger.Debug().Msg("Control point removed")
		ps.etcdClient.Delete(ps.statusEtcdPath)
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
			PolicyName:  ps.GetPolicyName(),
			PolicyHash:  ps.GetPolicyHash(),
			ComponentId: ps.GetComponentId(),
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
	ps.etcdClient.Put(ps.statusEtcdPath, string(data))
}
