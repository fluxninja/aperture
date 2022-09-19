package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiWatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/watch"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/entitycache/v1"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/utils"
)

const podTrackerPrefix = "kubernetes_pod"

type serviceCacheOperation int

const (
	add serviceCacheOperation = iota
	remove
)

type podInfo struct {
	Name      string
	Namespace string
	UID       string
	IPAddress string
}

func podInfoOrder(a, b podInfo) bool {
	if a.Namespace == b.Namespace {
		return a.Name < b.Name
	}
	return a.Namespace < b.Namespace
}

type serviceData struct {
	Name string
	Pods []podInfo
}

// serviceDataCache maps service UID to its data - the name of the service and pods handled by the service.
type serviceDataCache struct {
	cache map[string]*serviceData
}

func newServiceDataCache() *serviceDataCache {
	return &serviceDataCache{
		cache: make(map[string]*serviceData),
	}
}

func createServiceData(service *v1.Endpoints, nodeName string) *serviceData {
	fqdn := getFQDN(service)
	pods := getServicePods(service, nodeName)
	sort.Slice(pods, func(i, j int) bool {
		return podInfoOrder(pods[i], pods[j])
	})

	return &serviceData{
		Name: fqdn,
		Pods: pods,
	}
}

func (sc *serviceDataCache) updateService(service *v1.Endpoints, nodeName string) {
	uid := string(service.UID)
	serviceData := createServiceData(service, nodeName)

	sc.cache[uid] = serviceData
}

func (sc *serviceDataCache) updateServicePods(service *v1.Endpoints, pods []podInfo) {
	uid := string(service.UID)
	fqdn := getFQDN(service)
	sc.cache[uid] = &serviceData{
		Name: fqdn,
		Pods: pods,
	}
}

func (sc *serviceDataCache) removeService(service *v1.Endpoints) {
	uid := string(service.UID)
	delete(sc.cache, uid)
}

func (sc *serviceDataCache) getServiceData(service *v1.Endpoints) (*serviceData, bool) {
	uid := string(service.UID)
	val, ok := sc.cache[uid]
	return val, ok
}

// servicePodMapping maps pod (designated by namespace and pod name) to the list of its services.
type servicePodMapping struct {
	mapping map[string]map[string][]string
}

func newServicePodMapping() *servicePodMapping {
	return &servicePodMapping{
		mapping: make(map[string]map[string][]string),
	}
}

func (m *servicePodMapping) getServices(namespace, podName string) []string {
	return m.mapping[namespace][podName]
}

func (m *servicePodMapping) addService(namespace, podName, service string) {
	if _, ok := m.mapping[namespace]; !ok {
		m.mapping[namespace] = make(map[string][]string)
	}
	if _, ok := m.mapping[namespace][podName]; !ok {
		m.mapping[namespace][podName] = nil
	}
	// add to the list of services if it doesn't exist
	if !utils.SliceContains(m.mapping[namespace][podName], service) {
		m.mapping[namespace][podName] = append(m.mapping[namespace][podName], service)
	}
}

func (m *servicePodMapping) removeService(namespace, podName, service string) {
	if _, ok := m.mapping[namespace]; !ok {
		return
	}
	if _, ok := m.mapping[namespace][podName]; !ok {
		return
	}
	// remove from the list of services if it exists
	for i, s := range m.mapping[namespace][podName] {
		if s == service {
			m.mapping[namespace][podName] = append(m.mapping[namespace][podName][:i], m.mapping[namespace][podName][i+1:]...)
			break
		}
	}
	// if service list is empty, remove the pod
	if len(m.mapping[namespace][podName]) == 0 {
		delete(m.mapping[namespace], podName)
	}
	// if namespace is empty, remove the namespace
	if len(m.mapping[namespace]) == 0 {
		delete(m.mapping, namespace)
	}
}

// KubernetesDiscovery is a collector that collects Kubernetes information periodically.
type KubernetesDiscovery struct {
	waitGroup              sync.WaitGroup
	cli                    kubernetes.Interface
	ctx                    context.Context
	cancel                 context.CancelFunc
	trackers               notifiers.Trackers
	mapping                *servicePodMapping
	serviceCache           *serviceDataCache
	nodeName               string
	revisionEndpointsWatch string
}

func newKubernetesServiceDiscovery(trackers notifiers.Trackers, nodeName string, k8sClient k8s.K8sClient) (*KubernetesDiscovery, error) {
	if k8sClient.GetErrNotInCluster() {
		log.Info().Msg("Not in Kubernetes cluster, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}
	if k8sClient.GetErr() != nil {
		log.Error().Err(k8sClient.GetErr()).Msg("Error when creating Kubernetes client, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}

	if nodeName == "" {
		log.Error().Err(k8sClient.GetErr()).Msg("Node name not set, could not create Kubernetes service discovery")
		return nil, fmt.Errorf("node name not set")
	}

	kc := &KubernetesDiscovery{
		cli:          k8sClient.GetClientSet(),
		nodeName:     nodeName,
		mapping:      newServicePodMapping(),
		serviceCache: newServiceDataCache(),
		trackers:     trackers,
	}
	return kc, nil
}

func (kc *KubernetesDiscovery) start() {
	kc.ctx, kc.cancel = context.WithCancel(context.Background())

	kc.waitGroup.Add(1)

	panichandler.Go(func() {
		defer kc.waitGroup.Done()

		operation := func() error {
			// purge notifiers
			kc.trackers.Purge(podTrackerPrefix)

			// bootstrap mapping
			endpoints, err := kc.cli.CoreV1().Endpoints(metav1.NamespaceAll).List(kc.ctx, metav1.ListOptions{})
			if err != nil {
				log.Error().Err(err).Msg("Failed to list endpoints")
				return err
			}
			kc.revisionEndpointsWatch = endpoints.ResourceVersion
			if len(kc.mapping.mapping) > 0 {
				kc.mapping = newServicePodMapping()
			}
			for _, eItem := range endpoints.Items {
				e := eItem
				kc.addRemoveFromEndpoints(&e, add)
			}

			// setup watchers
			endpointsWatchFunc := func(options metav1.ListOptions) (apiWatch.Interface, error) {
				return kc.cli.CoreV1().Endpoints(metav1.NamespaceAll).Watch(kc.ctx, options)
			}
			var endpointsWatcher *watch.RetryWatcher

			endpointsWatcher, err = watch.NewRetryWatcher(kc.revisionEndpointsWatch,
				&cache.ListWatch{WatchFunc: endpointsWatchFunc})
			if err != nil {
				log.Error().Err(err).Msg("Unable to watch endpoints")
				return err
			}
			defer endpointsWatcher.Stop()

			for {
				// watchers added, start watching events
				select {
				case endpointEvent, ok := <-endpointsWatcher.ResultChan():
					if !ok {
						log.Error().Msg("Endpoints watcher closed")
						return fmt.Errorf("endpoints watcher closed")
					}
					switch endpointEvent.Type {
					case apiWatch.Added:
						endpoints := endpointEvent.Object.(*v1.Endpoints)
						kc.serviceCache.updateService(endpoints, kc.nodeName)
						kc.addRemoveFromEndpoints(endpoints, add)
					case apiWatch.Modified:
						endpoints := endpointEvent.Object.(*v1.Endpoints)
						cachedData, ok := kc.serviceCache.getServiceData(endpoints)
						if !ok {
							endpoints := endpointEvent.Object.(*v1.Endpoints)
							kc.serviceCache.updateService(endpoints, kc.nodeName)
							kc.addRemoveFromEndpoints(endpoints, add)
						} else {
							if getFQDN(endpoints) != cachedData.Name {
								kc.renameService(endpoints)
							}
							kc.syncPodLists(endpoints)
						}
					case apiWatch.Deleted:
						endpoints := endpointEvent.Object.(*v1.Endpoints)
						kc.addRemoveFromEndpoints(endpoints, remove)
						kc.serviceCache.removeService(endpoints)
					case apiWatch.Error:
						log.Error().Msg("Endpoints watcher error")
						return fmt.Errorf("endpoints watcher error")
					}
				case <-kc.ctx.Done():
					log.Info().Msg("KubeClient stopped")
					return backoff.Permanent(nil)
				}
			}
		}
		boff := backoff.NewConstantBackOff(5 * time.Second)

		_ = backoff.Retry(operation, backoff.WithContext(boff, kc.ctx))
		log.Info().Msg("Stopping kubernetes watcher")
	})
}

func (kc *KubernetesDiscovery) stop() {
	kc.cancel()
	kc.waitGroup.Wait()
	kc.trackers.Purge(podTrackerPrefix)
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kc *KubernetesDiscovery) updatePodInTracker(podInfo podInfo) error {
	services := kc.mapping.getServices(podInfo.Namespace, podInfo.Name)
	key := notifiers.Key(getPodIDKey(podInfo.UID))
	currentPodData := kc.trackers.GetCurrentValue(key)

	var entity *entitycachev1.Entity

	if len(currentPodData) > 0 {
		err := json.Unmarshal(currentPodData, &entity)
		if err != nil {
			log.Error().Msgf("Error unmarshalling entity: %v", err)
			return err
		}
	} else {
		entity = &entitycachev1.Entity{}
	}

	entity.Services = services
	entity.IpAddress = podInfo.IPAddress
	entity.Uid = podInfo.UID
	entity.Prefix = podTrackerPrefix
	entity.Name = podInfo.Name

	value, err := json.Marshal(entity)
	if err != nil {
		log.Error().Msgf("Error marshaling entity: %v", err)
		return err
	}
	kc.trackers.WriteEvent(key, value)
	return nil
}

func (kc *KubernetesDiscovery) syncPodLists(e *v1.Endpoints) {
	serviceName := getFQDN(e)
	cachedService, _ := kc.serviceCache.getServiceData(e)

	// assume cached pods are sorted by namespace and name
	cachedPods := cachedService.Pods

	currentPods := getServicePods(e, kc.nodeName)
	sort.Slice(currentPods, func(i, j int) bool {
		return podInfoOrder(currentPods[i], currentPods[j])
	})

	cacheIndexOffset := 0

	for i := 0; i < len(currentPods); i++ {
		cachePodIndex := i + cacheIndexOffset
		currentPod := currentPods[i]
		var cachePod podInfo
		outOfCache := false
		if cachePodIndex < len(cachedPods) {
			cachePod = cachedPods[cachePodIndex]
		} else {
			outOfCache = true
		}
		if currentPod == cachePod {
			continue
		}
		if outOfCache || podInfoOrder(currentPod, cachePod) {
			// a pod is missing from cachedPods slice - it was added to the service
			kc.mapping.addService(currentPod.Namespace, currentPod.Name, serviceName)
			cacheIndexOffset--

			err := kc.updatePodInTracker(currentPod)
			if err != nil {
				log.Error().Msgf("Tracker could not be updated: %v", err)
			}
		} else {
			// a pod is missing from currentPods slice - it should be removed from the service
			kc.mapping.removeService(cachePod.Namespace, cachePod.Name, serviceName)
			i--
			cacheIndexOffset++

			err := kc.updatePodInTracker(cachePod)
			if err != nil {
				log.Error().Msgf("Tracker could not be updated: %v", err)
			}
		}
	}

	for i := len(currentPods) + cacheIndexOffset; i < len(cachedPods); i++ {
		cachePod := cachedPods[i]
		// a pod is missing from currentPods slice - it should be removed from the service
		kc.mapping.removeService(cachePod.Namespace, cachePod.Name, serviceName)

		err := kc.updatePodInTracker(cachePod)
		if err != nil {
			log.Error().Msgf("Tracker could not be updated: %v", err)
		}
	}

	kc.serviceCache.updateServicePods(e, currentPods)
}

// getFQDN return the full qualified domain name of a given service.
func getFQDN(e *v1.Endpoints) string {
	name := e.Name
	namespace := e.Namespace

	// we assume that FQDN of all kubernetes services is the default one
	defaultFQDN := fmt.Sprintf("%v.%v.svc.cluster.local", name, namespace)
	return defaultFQDN
}

// getServicePods retrieves a list of pods handled by a given service that are located on a given node.
func getServicePods(service *v1.Endpoints, nodeName string) []podInfo {
	var pods []podInfo

	for _, subset := range service.Subsets {
		for _, address := range subset.Addresses {
			if address.TargetRef == nil {
				continue
			}
			if address.TargetRef.Kind != "Pod" {
				continue
			}
			if *(address.NodeName) != nodeName {
				continue
			}
			p := podInfo{
				Name:      address.TargetRef.Name,
				Namespace: address.TargetRef.Namespace,
				UID:       string(address.TargetRef.UID),
				IPAddress: address.IP,
			}
			pods = append(pods, p)
		}
	}

	return pods
}

func (kc *KubernetesDiscovery) renameService(e *v1.Endpoints) {
	newServiceName := getFQDN(e)
	oldService, _ := kc.serviceCache.getServiceData(e)
	oldServiceName := oldService.Name
	for _, podInfo := range oldService.Pods {
		kc.mapping.removeService(podInfo.Namespace, podInfo.Name, oldServiceName)
		kc.mapping.addService(podInfo.Namespace, podInfo.Name, newServiceName)

		err := kc.updatePodInTracker(podInfo)
		if err != nil {
			log.Error().Msgf("Tracker could not be updated: %v", err)
		}
	}
}

func (kc *KubernetesDiscovery) addRemoveFromEndpoints(e *v1.Endpoints, operation serviceCacheOperation) {
	serviceName := getFQDN(e)
	pods := getServicePods(e, kc.nodeName)
	for _, pod := range pods {
		if operation == add {
			kc.mapping.addService(pod.Namespace, pod.Name, serviceName)
		} else {
			kc.mapping.removeService(pod.Namespace, pod.Name, serviceName)
		}
		err := kc.updatePodInTracker(pod)
		if err != nil {
			log.Error().Msgf("Tracker could not be updated: %v", err)
		}
	}
}

func getPodIDKey(key string) string {
	return fmt.Sprintf("%s.%s", podTrackerPrefix, key)
}
