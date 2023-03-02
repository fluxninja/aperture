package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/cenkalti/backoff/v4"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/entitycache/v1"
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
	Node      string
	UID       string
	IPAddress string
	ClusterIP string
}

func podInfoOrder(a, b podInfo) bool {
	if a.Namespace == b.Namespace {
		return a.Name < b.Name
	}
	return a.Namespace < b.Namespace
}

type serviceData struct {
	fqdn string
	pods []podInfo
}

// serviceDataCache maps service UID to its data - the name of the service and pods handled by the service.
type serviceDataCache struct {
	cache map[string][]podInfo
	kd    *serviceDiscovery
}

func newServiceDataCache(kd *serviceDiscovery) *serviceDataCache {
	return &serviceDataCache{
		cache: make(map[string][]podInfo),
		kd:    kd,
	}
}

func createServiceData(kd *serviceDiscovery, endpoints *v1.Endpoints) *serviceData {
	fqdn := kd.getFQDN(endpoints)
	pods := getServicePods(endpoints)
	sort.Slice(pods, func(i, j int) bool {
		return podInfoOrder(pods[i], pods[j])
	})

	return &serviceData{
		fqdn: fqdn,
		pods: pods,
	}
}

func (sc *serviceDataCache) updateServiceData(endpoints *v1.Endpoints) {
	serviceData := createServiceData(sc.kd, endpoints)
	sc.cache[serviceData.fqdn] = serviceData.pods
}

func (sc *serviceDataCache) removeService(endpoints *v1.Endpoints) {
	fqdn := sc.kd.getFQDN(endpoints)
	delete(sc.cache, fqdn)
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

func (m *servicePodMapping) getFQDNs(namespace, podName string) []string {
	return m.mapping[namespace][podName]
}

func (m *servicePodMapping) addService(namespace, podName, fqdn string) {
	if _, ok := m.mapping[namespace]; !ok {
		m.mapping[namespace] = make(map[string][]string)
	}
	if _, ok := m.mapping[namespace][podName]; !ok {
		m.mapping[namespace][podName] = nil
	}
	// add to the list of services if it doesn't exist
	if !utils.SliceContains(m.mapping[namespace][podName], fqdn) {
		m.mapping[namespace][podName] = append(m.mapping[namespace][podName], fqdn)
	}
}

func (m *servicePodMapping) removeService(namespace, podName, fqdn string) {
	if _, ok := m.mapping[namespace]; !ok {
		return
	}
	if _, ok := m.mapping[namespace][podName]; !ok {
		return
	}
	// remove from the list of services if it exists
	for i, s := range m.mapping[namespace][podName] {
		if s == fqdn {
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

// serviceDiscovery is a collector that collects Kubernetes information periodically.
type serviceDiscovery struct {
	clusterDomain     string
	cli               kubernetes.Interface
	informerFactory   informers.SharedInformerFactory
	informersStopChan chan struct{}
	entityEvents      notifiers.EventWriter
	mapping           *servicePodMapping
	serviceCache      *serviceDataCache
}

func newServiceDiscovery(
	entityEvents notifiers.EventWriter,
	k8sClient k8s.K8sClient,
) (*serviceDiscovery, error) {
	kd := &serviceDiscovery{
		cli:          k8sClient.GetClientSet(),
		mapping:      newServicePodMapping(),
		entityEvents: entityEvents,
	}
	kd.serviceCache = newServiceDataCache(kd)
	kd.informerFactory = informers.NewSharedInformerFactory(kd.cli, 0)
	kd.informersStopChan = make(chan struct{})
	return kd, nil
}

func (kd *serviceDiscovery) start(ctx context.Context) {
	panichandler.Go(func() {
		operation := func() error {
			// get cluster domain
			clusterDomain, err := utils.GetClusterDomain()
			if err != nil {
				log.Error().Err(err).Msg("Could not get cluster domain, will retry")
				return err
			}
			kd.clusterDomain = clusterDomain

			// purge notifiers
			kd.entityEvents.Purge("")

			endpointsInformer := kd.informerFactory.Core().V1().Endpoints().Informer()
			_, err = endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc:    kd.handleEndpointsAdd,
				UpdateFunc: kd.handleEndpointsUpdate,
				DeleteFunc: kd.handleEndpointsDelete,
			})
			if err != nil {
				return err
			}

			kd.informerFactory.Start(kd.informersStopChan)
			if !cache.WaitForCacheSync(kd.informersStopChan, endpointsInformer.HasSynced) {
				return fmt.Errorf("timed out waiting for caches to sync")
			}

			<-kd.informersStopChan
			return nil
		}
		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, ctx))

		log.Info().Msg("Stopping kubernetes service watcher")
	})
}

func (kd *serviceDiscovery) handleEndpointsAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.serviceCache.updateServiceData(endpoints)
	kd.updateMappingFromEndpoints(endpoints, add)
}

func (kd *serviceDiscovery) handleEndpointsUpdate(oldObj, newObj interface{}) {
	oldEndpoints := oldObj.(*v1.Endpoints)
	newEndpoints := newObj.(*v1.Endpoints)
	if oldEndpoints.ResourceVersion == newEndpoints.ResourceVersion {
		return
	}
	fqdn := kd.getFQDN(oldEndpoints)
	cachedPods, ok := kd.serviceCache.cache[fqdn]
	if ok {
		for _, pod := range cachedPods {
			kd.mapping.removeService(pod.Namespace, pod.Name, fqdn)
			kd.removeEntityFromTracker(pod)
		}
		kd.serviceCache.removeService(oldEndpoints)
	}
	kd.serviceCache.updateServiceData(newEndpoints)
	kd.updateMappingFromEndpoints(newEndpoints, add)
}

func (kd *serviceDiscovery) handleEndpointsDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.updateMappingFromEndpoints(endpoints, remove)
	kd.serviceCache.removeService(endpoints)
}

func (kd *serviceDiscovery) stop(ctx context.Context) {
	close(kd.informersStopChan)
	ctx.Done()
	kd.entityEvents.Purge("")
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) writeEntityInTracker(podInfo podInfo) error {
	services := kd.mapping.getFQDNs(podInfo.Namespace, podInfo.Name)
	entity := &entitycachev1.Entity{
		Services:  services,
		IpAddress: podInfo.IPAddress,
		Uid:       podInfo.UID,
		Prefix:    podTrackerPrefix,
		Name:      podInfo.Name,
		ClusterIp: podInfo.ClusterIP,
	}

	value, err := json.Marshal(entity)
	if err != nil {
		log.Error().Msgf("Error marshaling entity: %v", err)
		return err
	}

	kd.entityEvents.WriteEvent(notifiers.Key(podInfo.UID), value)
	return nil
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) removeEntityFromTracker(podInfo podInfo) {
	kd.entityEvents.RemoveEvent(notifiers.Key(podInfo.UID))
}

func (kd *serviceDiscovery) updateMappingFromEndpoints(endpoints *v1.Endpoints, operation serviceCacheOperation) {
	fqdn := kd.getFQDN(endpoints)
	pods := getServicePods(endpoints)

	for _, pod := range pods {
		if operation == add {
			kd.mapping.addService(pod.Namespace, pod.Name, fqdn)
			err := kd.writeEntityInTracker(pod)
			if err != nil {
				log.Error().Msgf("Tracker could not be updated: %v", err)
			}
		} else {
			kd.mapping.removeService(pod.Namespace, pod.Name, fqdn)
			kd.removeEntityFromTracker(pod)
		}
	}
}

// getFQDN return the full qualified domain name of a given service.
func (kd *serviceDiscovery) getFQDN(endpoints *v1.Endpoints) string {
	name := endpoints.Name
	namespace := endpoints.Namespace

	serviceFQDN := fmt.Sprintf("%s.%s.svc.%s", name, namespace, kd.clusterDomain)
	return serviceFQDN
}

// getServicePods retrieves a list of pods handled by a given service that are located on a given node.
func getServicePods(endpoints *v1.Endpoints) []podInfo {
	var pods []podInfo

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			if address.TargetRef == nil {
				continue
			}
			if address.TargetRef.Kind != "Pod" {
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
