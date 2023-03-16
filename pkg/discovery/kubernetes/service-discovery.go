package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sourcegraph/conc/stream"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/utils"
)

const podTrackerPrefix = "kubernetes_pod"

// serviceDiscovery is a collector that collects Kubernetes information periodically.
type serviceDiscovery struct {
	ctx             context.Context
	cli             kubernetes.Interface
	informerFactory informers.SharedInformerFactory
	entityEvents    notifiers.EventWriter
	cancel          context.CancelFunc
	serviceStream   *stream.Stream
	clusterDomain   string

	podCounter *prometheus.GaugeVec
	election   *election.Election
}

func newServiceDiscovery(
	entityEvents notifiers.EventWriter,
	k8sClient k8s.K8sClient,
	pr *prometheus.Registry,
	election *election.Election,
) (*serviceDiscovery, error) {
	kd := &serviceDiscovery{
		cli:           k8sClient.GetClientSet(),
		entityEvents:  entityEvents,
		serviceStream: stream.New(),
		election:      election,
	}
	if kd.election != nil && kd.election.IsLeader() {
		defaultLabels := []string{
			metrics.K8sNamespaceName, metrics.K8sNodeName, metrics.K8sStatus,
			metrics.K8sCronjobName, metrics.K8sDaemonsetName, metrics.K8sDeploymentName, metrics.K8sJobName, metrics.K8sReplicasetName, metrics.K8sStatefulsetName,
		}
		podCounter := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metrics.K8sPodCount,
			Help: "The number of pods",
		}, defaultLabels)
		err := pr.Register(podCounter)
		if err != nil {
			// Ignore already registered error, as this is not harmful. Metrics may
			// be registered by other running server.
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, fmt.Errorf("couldn't register prometheus metrics: %w", err)
			}
		}
		kd.podCounter = podCounter
	}
	kd.informerFactory = informers.NewSharedInformerFactory(kd.cli, 0)
	kd.ctx, kd.cancel = context.WithCancel(context.Background())
	return kd, nil
}

func (kd *serviceDiscovery) start(startCtx context.Context) error {
	// get cluster domain
	clusterDomain, err := utils.GetClusterDomain()
	if err != nil {
		log.Error().Err(err).Msg("Could not get cluster domain, will retry")
		return err
	}
	kd.clusterDomain = clusterDomain
	var informerSynced []cache.InformerSynced
	endpointsInformer := kd.informerFactory.Core().V1().Endpoints().Informer()
	_, err = endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    kd.handleEndpointsAdd,
		UpdateFunc: kd.handleEndpointsUpdate,
		DeleteFunc: kd.handleEndpointsDelete,
	})
	if err != nil {
		return err
	}
	informerSynced = append(informerSynced, endpointsInformer.HasSynced)
	serviceInformer := kd.informerFactory.Core().V1().Services().Informer()
	_, err = serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    kd.handleServiceAdd,
		DeleteFunc: kd.handleServiceDelete,
	})
	if err != nil {
		return err
	}
	informerSynced = append(informerSynced, serviceInformer.HasSynced)
	if kd.election != nil && kd.election.IsLeader() {
		podInformer := kd.informerFactory.Core().V1().Pods().Informer()
		_, err = podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc:    kd.handlePodAdd,
			UpdateFunc: kd.handlePodUpdate,
			DeleteFunc: kd.handlePodDelete,
		})
		if err != nil {
			return err
		}
		informerSynced = append(informerSynced, podInformer.HasSynced)
	}
	kd.informerFactory.Start(kd.ctx.Done())

	if !cache.WaitForCacheSync(startCtx.Done(), informerSynced...) {
		return errors.New("timed out waiting for caches to sync")
	}

	log.Info().Msg("Service discovery started")
	return nil
}

func (kd *serviceDiscovery) stop(stopCtx context.Context) error {
	kd.cancel()
	kd.serviceStream.Wait()
	kd.entityEvents.Purge("")
	return nil
}

// Endpoints informer handlers

func (kd *serviceDiscovery) handleEndpointsAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.updateEndpoints(endpoints)
}

func (kd *serviceDiscovery) handleEndpointsUpdate(oldObj, newObj interface{}) {
	oldEndpoints := oldObj.(*v1.Endpoints)
	newEndpoints := newObj.(*v1.Endpoints)

	// remove addresses that are no longer in the newEndpoints
	for _, oldSubset := range oldEndpoints.Subsets {
		for _, oldAddress := range oldSubset.Addresses {
			found := false
			for _, newSubset := range newEndpoints.Subsets {
				for _, newAddress := range newSubset.Addresses {
					if newAddress.TargetRef.UID == oldAddress.TargetRef.UID {
						found = true
						break
					}
				}
			}
			if !found {
				// address is no longer in the newEndpoints, remove it
				kd.removeEndpointAddress(oldAddress, oldEndpoints.Namespace, oldEndpoints.Name)
			}
		}
	}

	kd.updateEndpoints(newEndpoints)
}

func (kd *serviceDiscovery) handleEndpointsDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.removeEndpoints(endpoints)
}

func (kd *serviceDiscovery) updateEndpoints(endpoints *v1.Endpoints) {
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			if address.TargetRef == nil {
				continue
			}
			if address.TargetRef.Kind != "Pod" {
				continue
			}
			p := &entitiesv1.Entity{
				Name:      address.TargetRef.Name,
				Namespace: address.TargetRef.Namespace,
				Uid:       string(address.TargetRef.UID),
				IpAddress: address.IP,
				Services:  []string{kd.getService(endpoints.Namespace, endpoints.Name)},
			}
			if address.NodeName != nil {
				p.NodeName = *address.NodeName
			}
			err := kd.updateEntity(p)
			if err != nil {
				log.Error().Err(err).Msg("Tracker could not be updated")
			}
		}
	}
}

// Service informer handlers

func (kd *serviceDiscovery) handleServiceAdd(obj interface{}) {
	service := obj.(*v1.Service)
	kd.addClusterIPs(service)
}

func (kd *serviceDiscovery) handleServiceDelete(obj interface{}) {
	service := obj.(*v1.Service)
	kd.removeClusterIPs(service)
}

func (kd *serviceDiscovery) addClusterIPs(service *v1.Service) {
	clusterIPs := getClusterIPsFromService(service)
	for _, clusterIP := range clusterIPs {
		p := &entitiesv1.Entity{
			Name:      strings.Join([]string{"ClusterIP", service.Namespace, service.Name, clusterIP}, "-"),
			Namespace: service.Namespace,
			Uid:       clusterIP,
			IpAddress: clusterIP,
			Services:  []string{kd.getService(service.Namespace, service.Name)},
		}
		err := kd.updateEntity(p)
		if err != nil {
			log.Error().Err(err).Msg("Tracker could not be updated")
		}
	}
}

func (kd *serviceDiscovery) removeClusterIPs(service *v1.Service) {
	clusterIPs := getClusterIPsFromService(service)
	for _, clusterIP := range clusterIPs {
		kd.entityEvents.RemoveEvent(notifiers.Key(clusterIP))
	}
}

// getService return the full qualified domain name of a given service.
func (kd *serviceDiscovery) getService(namespace, name string) string {
	service := fmt.Sprintf("%s.%s.svc.%s", name, namespace, kd.clusterDomain)
	return service
}

func (kd *serviceDiscovery) removeEndpointAddress(address v1.EndpointAddress, namespace, name string) {
	updateFunc := func(oldValue []byte) (notifiers.EventType, []byte) {
		if oldValue == nil {
			return notifiers.Remove, nil
		}
		entity := &entitiesv1.Entity{}
		err := json.Unmarshal(oldValue, entity)
		if err != nil {
			log.Error().Err(err).Msg("Could not unmarshal entity")
			return notifiers.Remove, nil
		}
		entity.Services = utils.RemoveFromSlice(entity.Services, kd.getService(namespace, name))
		if shouldRemove(entity) {
			return notifiers.Remove, nil
		}
		bytes, err := json.Marshal(entity)
		if err != nil {
			log.Error().Err(err).Msg("Could not marshal entity")
			return notifiers.Remove, nil
		}
		return notifiers.Write, bytes
	}

	kd.entityEvents.UpdateValue(notifiers.Key(address.TargetRef.UID), updateFunc)
}

func (kd *serviceDiscovery) removeEndpoints(endpoints *v1.Endpoints) {
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			kd.removeEndpointAddress(address, endpoints.Namespace, endpoints.Name)
		}
	}
}

// updateEntity updates the entity in the tracker.
func (kd *serviceDiscovery) updateEntity(pod *entitiesv1.Entity) error {
	updateFunc := func(oldValue []byte) (notifiers.EventType, []byte) {
		var entity *entitiesv1.Entity
		if oldValue == nil {
			// create new entity
			entity = pod
		} else {
			err := json.Unmarshal(oldValue, &entity)
			if err != nil {
				log.Error().Msgf("Error unmarshaling entity: %v", err)
				return notifiers.Write, oldValue
			}
			for _, service := range pod.Services {
				if !utils.SliceContains(entity.Services, service) {
					entity.Services = append(entity.Services, service)
				}
			}
		}
		value, err := json.Marshal(entity)
		if err != nil {
			log.Error().Msgf("Error marshaling entity: %v", err)
			return notifiers.Write, oldValue
		}
		return notifiers.Write, value
	}

	kd.entityEvents.UpdateValue(notifiers.Key(pod.Uid), updateFunc)

	return nil
}

func (kd *serviceDiscovery) handlePodAdd(obj any) {
	pod := obj.(*v1.Pod)
	labels := podLabels(pod)

	podCounter, err := kd.podCounter.GetMetricWith(labels)
	if err != nil {
		log.Debug().Msgf("Could not extract request counter metric from registry: %v", err)
		return
	}
	podCounter.Inc()
}

func (kd *serviceDiscovery) handlePodDelete(obj any) {
	pod := obj.(*v1.Pod)
	newPod := pod.DeepCopy()
	newPod.Status.Phase = v1.PodSucceeded
	kd.handlePodUpdate(pod, newPod)
}

func (kd *serviceDiscovery) handlePodUpdate(oldObj, newObj any) {
	oldPod := oldObj.(*v1.Pod)
	newPod := newObj.(*v1.Pod)

	if oldPod.Status.Phase != newPod.Status.Phase || oldPod.Spec.NodeName != newPod.Spec.NodeName {
		oldPodCounter, errOld := kd.podCounter.GetMetricWith(podLabels(oldPod))
		newPodCounter, errNew := kd.podCounter.GetMetricWith(podLabels(newPod))
		if errOld != nil {
			log.Error().Msgf("Could not extract request counter metric from registry: %v", errOld)
			return
		}
		if errNew != nil {
			log.Error().Msgf("Could not extract request counter metric from registry: %v", errNew)
			return
		}
		oldPodCounter.Dec()
		newPodCounter.Inc()
	}
}

/* Helper functions */

func podLabels(pod *v1.Pod) map[string]string {
	labels := map[string]string{
		metrics.K8sNamespaceName:   pod.Namespace,
		metrics.K8sNodeName:        pod.Spec.NodeName,
		metrics.K8sStatus:          string(pod.Status.Phase),
		metrics.K8sCronjobName:     "",
		metrics.K8sDaemonsetName:   "",
		metrics.K8sDeploymentName:  "",
		metrics.K8sJobName:         "",
		metrics.K8sReplicasetName:  "",
		metrics.K8sStatefulsetName: "",
	}
	owners := pod.GetObjectMeta().GetOwnerReferences()
	for _, owner := range owners {
		kind := strings.ToLower(owner.Kind)
		labels[fmt.Sprintf("k8s_%s_name", kind)] = owner.Name
	}
	return labels
}

func getClusterIPsFromService(svc *v1.Service) []string {
	if svc.Spec.Type != v1.ServiceTypeClusterIP {
		return nil
	}
	clusterIPs := svc.Spec.ClusterIPs
	// remove None, "" from the list
	clusterIPs = utils.RemoveFromSlice(clusterIPs, "None")
	clusterIPs = utils.RemoveFromSlice(clusterIPs, "")
	return clusterIPs
}

func shouldRemove(entity *entitiesv1.Entity) bool {
	// once we have more informers then add additional checks here
	return len(entity.Services) == 0
}
