package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sourcegraph/conc/stream"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	UID       string
	IPAddress string
	ClusterIP string
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
	ctx             context.Context
	cancel          context.CancelFunc
	clusterDomain   string
	cli             kubernetes.Interface
	informerFactory informers.SharedInformerFactory
	entityEvents    notifiers.EventWriter
	mapping         *servicePodMapping
	serviceStream   *stream.Stream
}

func newServiceDiscovery(
	entityEvents notifiers.EventWriter,
	k8sClient k8s.K8sClient,
) (*serviceDiscovery, error) {
	kd := &serviceDiscovery{
		cli:           k8sClient.GetClientSet(),
		mapping:       newServicePodMapping(),
		entityEvents:  entityEvents,
		serviceStream: stream.New(),
	}
	kd.informerFactory = informers.NewSharedInformerFactory(kd.cli, 0)
	kd.ctx, kd.cancel = context.WithCancel(context.Background())
	return kd, nil
}

func (kd *serviceDiscovery) start() {
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

			kd.informerFactory.Start(kd.ctx.Done())
			if !cache.WaitForCacheSync(kd.ctx.Done(), endpointsInformer.HasSynced) {
				return fmt.Errorf("timed out waiting for caches to sync")
			}

			<-kd.ctx.Done()
			return nil
		}
		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, kd.ctx))
		log.Info().Msg("Service discovery stopped")
	})
}

func (kd *serviceDiscovery) stop() {
	kd.cancel()
	kd.serviceStream.Wait()
	kd.entityEvents.Purge("")
}

func (kd *serviceDiscovery) handleEndpointsAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.updateMappingFromEndpoints(endpoints, add)
}

func (kd *serviceDiscovery) handleEndpointsUpdate(oldObj, newObj interface{}) {
	oldEndpoints := oldObj.(*v1.Endpoints)
	newEndpoints := newObj.(*v1.Endpoints)
	if oldEndpoints.ResourceVersion == newEndpoints.ResourceVersion {
		return
	}

	// remove old services
	fqdn := getFQDN(oldEndpoints, kd.clusterDomain)
	for _, oldSubsets := range oldEndpoints.Subsets {
		for _, oldAddress := range oldSubsets.Addresses {
			kd.mapping.removeService(oldAddress.TargetRef.Namespace, oldAddress.TargetRef.Name, fqdn)
			kd.removeEntityFromTracker(string(oldAddress.TargetRef.UID))
		}
	}

	// add new services
	kd.updateMappingFromEndpoints(newEndpoints, add)
}

func (kd *serviceDiscovery) handleEndpointsDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.updateMappingFromEndpoints(endpoints, remove)
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) writeEntityInTracker(pod podInfo) error {
	fqdns := kd.mapping.getFQDNs(pod.Namespace, pod.Name)
	entity := &entitycachev1.Entity{
		Services:  fqdns,
		IpAddress: pod.IPAddress,
		Uid:       pod.UID,
		Prefix:    podTrackerPrefix,
		Name:      pod.Name,
		Labels: map[string]string{
			"cluster_ip": pod.ClusterIP,
		},
	}

	value, err := json.Marshal(entity)
	if err != nil {
		log.Error().Msgf("Error marshaling entity: %v", err)
		return err
	}

	kd.entityEvents.WriteEvent(notifiers.Key(pod.UID), value)
	return nil
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) removeEntityFromTracker(uid string) {
	kd.entityEvents.RemoveEvent(notifiers.Key(uid))
}

func (kd *serviceDiscovery) updateMappingFromEndpoints(endpoints *v1.Endpoints, operation serviceCacheOperation) {
	fqdn := getFQDN(endpoints, kd.clusterDomain)
	pods := getServicePods(endpoints)

	kd.serviceStream.Go(func() stream.Callback {
		return kd.fetchService(endpoints, pods, fqdn, operation)
	})
}

func (kd *serviceDiscovery) fetchService(endpoints *v1.Endpoints, pods []podInfo, fqdn string, operation serviceCacheOperation) stream.Callback {
	clusterIP := ""
	op := func() error {
		svc, err := kd.cli.CoreV1().Services(endpoints.Namespace).Get(kd.ctx, endpoints.Name, metav1.GetOptions{})
		if kd.ctx.Err() != nil {
			return backoff.Permanent(kd.ctx.Err())
		}
		if err != nil {
			log.Trace().Err(err).Str("namespace", endpoints.Namespace).Str("name", endpoints.Name).Msg("Could not fetch service")
			return nil
		}
		if svc.Spec.Type != v1.ServiceTypeClusterIP {
			return nil
		}
		clusterIP = svc.Spec.ClusterIP
		return nil
	}
	err := backoff.Retry(op, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
	if err != nil {
		log.Error().Err(err).Str("namespace", endpoints.Namespace).Str("name", endpoints.Name).Msg("Context canceled while fetching service")
		return func() {}
	}

	return func() {
		for _, pod := range pods {
			p := pod
			if clusterIP != "" && clusterIP != "None" {
				p.ClusterIP = clusterIP
			} else {
				p.ClusterIP = ""
			}
			if operation == add {
				kd.mapping.addService(p.Namespace, p.Name, fqdn)
				err := kd.writeEntityInTracker(p)
				if err != nil {
					log.Error().Err(err).Msg("Tracker could not be updated")
				}
			} else {
				kd.mapping.removeService(p.Namespace, p.Name, fqdn)
				kd.removeEntityFromTracker(p.UID)
			}
		}
	}
}

// getFQDN return the full qualified domain name of a given service.
func getFQDN(endpoints *v1.Endpoints, clusterDomain string) string {
	name := endpoints.Name
	namespace := endpoints.Namespace

	serviceFQDN := fmt.Sprintf("%s.%s.svc.%s", name, namespace, clusterDomain)
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
