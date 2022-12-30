package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiWatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/watch"

	entitycachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/entitycache/v1"
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
	kd    *KubernetesDiscovery
}

func newServiceDataCache(kd *KubernetesDiscovery) *serviceDataCache {
	return &serviceDataCache{
		cache: make(map[string][]podInfo),
		kd:    kd,
	}
}

func createServiceData(kd *KubernetesDiscovery, endpoints *v1.Endpoints, nodeName string) *serviceData {
	fqdn := kd.getFQDN(endpoints)
	pods := getServicePods(endpoints, nodeName)
	sort.Slice(pods, func(i, j int) bool {
		return podInfoOrder(pods[i], pods[j])
	})

	return &serviceData{
		fqdn: fqdn,
		pods: pods,
	}
}

func (sc *serviceDataCache) updateServiceData(endpoints *v1.Endpoints, nodeName string) {
	serviceData := createServiceData(sc.kd, endpoints, nodeName)
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

// KubernetesDiscovery is a collector that collects Kubernetes information periodically.
type KubernetesDiscovery struct {
	waitGroup              sync.WaitGroup
	cli                    kubernetes.Interface
	ctx                    context.Context
	cancel                 context.CancelFunc
	entityEvents           notifiers.EventWriter
	mapping                *servicePodMapping
	serviceCache           *serviceDataCache
	nodeName               string
	revisionEndpointsWatch string
	clusterDomain          string
}

func newKubernetesServiceDiscovery(
	entityEvents notifiers.EventWriter,
	nodeName string,
	k8sClient k8s.K8sClient,
) (*KubernetesDiscovery, error) {
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

	kd := &KubernetesDiscovery{
		cli:          k8sClient.GetClientSet(),
		nodeName:     nodeName,
		mapping:      newServicePodMapping(),
		entityEvents: entityEvents,
	}
	kd.serviceCache = newServiceDataCache(kd)
	return kd, nil
}

func (kd *KubernetesDiscovery) start() {
	kd.ctx, kd.cancel = context.WithCancel(context.Background())

	kd.waitGroup.Add(1)

	panichandler.Go(func() {
		defer kd.waitGroup.Done()

		operation := func() error {
			// get cluster domain
			clusterDomain, err := getClusterDomain()
			if err != nil {
				log.Error().Err(err).Msg("Could not get cluster domain, will retry")
				return err
			}
			kd.clusterDomain = clusterDomain

			// purge notifiers
			kd.entityEvents.Purge("")

			// bootstrap mapping
			endpoints, err := kd.cli.CoreV1().Endpoints(metav1.NamespaceAll).List(kd.ctx, metav1.ListOptions{})
			if err != nil {
				log.Error().Err(err).Msg("Failed to list endpoints")
				return err
			}
			kd.revisionEndpointsWatch = endpoints.ResourceVersion
			if len(kd.mapping.mapping) > 0 {
				kd.mapping = newServicePodMapping()
			}
			for _, item := range endpoints.Items {
				e := item
				kd.serviceCache.updateServiceData(&e, kd.nodeName)
				kd.updateMappingFromEndpoints(&e, add)
			}

			// setup watchers
			endpointsWatchFunc := func(options metav1.ListOptions) (apiWatch.Interface, error) {
				return kd.cli.CoreV1().Endpoints(metav1.NamespaceAll).Watch(kd.ctx, options)
			}
			endpointsWatcher, err := watch.NewRetryWatcher(kd.revisionEndpointsWatch, &cache.ListWatch{
				WatchFunc: endpointsWatchFunc,
			})
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
						kd.serviceCache.updateServiceData(endpoints, kd.nodeName)
						kd.updateMappingFromEndpoints(endpoints, add)
					case apiWatch.Modified:
						endpoints := endpointEvent.Object.(*v1.Endpoints)
						fqdn := kd.getFQDN(endpoints)
						cachedPods, ok := kd.serviceCache.cache[fqdn]
						if ok {
							for _, cachedPod := range cachedPods {
								kd.mapping.removeService(cachedPod.Namespace, cachedPod.Name, fqdn)
								kd.removeEntityFromTracker(cachedPod)
							}
							kd.serviceCache.removeService(endpoints)
						}
						kd.serviceCache.updateServiceData(endpoints, kd.nodeName)
						kd.updateMappingFromEndpoints(endpoints, add)
					case apiWatch.Deleted:
						endpoints := endpointEvent.Object.(*v1.Endpoints)
						kd.updateMappingFromEndpoints(endpoints, remove)
						kd.serviceCache.removeService(endpoints)
					case apiWatch.Error:
						log.Error().Msg("Endpoints watcher error")
						return fmt.Errorf("endpoints watcher error")
					}
				case <-kd.ctx.Done():
					log.Info().Msg("KubeClient stopped")
					return backoff.Permanent(nil)
				}
			}
		}
		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, kd.ctx))

		log.Info().Msg("Stopping kubernetes watcher")
	})
}

func (kd *KubernetesDiscovery) stop() {
	kd.cancel()
	kd.waitGroup.Wait()
	kd.entityEvents.Purge("")
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *KubernetesDiscovery) writeEntityInTracker(podInfo podInfo) error {
	services := kd.mapping.getFQDNs(podInfo.Namespace, podInfo.Name)
	entity := &entitycachev1.Entity{
		Services:  services,
		IpAddress: podInfo.IPAddress,
		Uid:       podInfo.UID,
		Prefix:    podTrackerPrefix,
		Name:      podInfo.Name,
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
func (kd *KubernetesDiscovery) removeEntityFromTracker(podInfo podInfo) {
	kd.entityEvents.RemoveEvent(notifiers.Key(podInfo.UID))
}

func (kd *KubernetesDiscovery) updateMappingFromEndpoints(endpoints *v1.Endpoints, operation serviceCacheOperation) {
	fqdn := kd.getFQDN(endpoints)
	pods := getServicePods(endpoints, kd.nodeName)

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
func (kd *KubernetesDiscovery) getFQDN(endpoints *v1.Endpoints) string {
	name := endpoints.Name
	namespace := endpoints.Namespace

	serviceFQDN := fmt.Sprintf("%s.%s.svc.%s", name, namespace, kd.clusterDomain)
	return serviceFQDN
}

// getServicePods retrieves a list of pods handled by a given service that are located on a given node.
func getServicePods(endpoints *v1.Endpoints, nodeName string) []podInfo {
	var pods []podInfo

	for _, subset := range endpoints.Subsets {
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
				Node:      *address.NodeName,
				UID:       string(address.TargetRef.UID),
				IPAddress: address.IP,
			}
			pods = append(pods, p)
		}
	}

	return pods
}

// Retrieve cluster domain of Kubernetes cluster we are installed on. It can be retrieved by looking up CNAME of
// kubernetes.default.svc and extracting its suffix.
func getClusterDomain() (string, error) {
	apiSvc := "kubernetes.default.svc"

	cname, err := net.LookupCNAME(apiSvc)
	if err != nil {
		return "", err
	}

	clusterDomain := strings.TrimPrefix(cname, apiSvc+".")
	clusterDomain = strings.TrimSuffix(clusterDomain, ".")

	return clusterDomain, nil
}
