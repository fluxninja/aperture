package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
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
	Services  []string
	ClusterIP string
}

// serviceDiscovery is a collector that collects Kubernetes information periodically.
type serviceDiscovery struct {
	waitGroup     sync.WaitGroup
	cli           kubernetes.Interface
	ctx           context.Context
	cancel        context.CancelFunc
	entityEvents  notifiers.EventWriter
	nodeName      string
	clusterDomain string
}

func newServiceDiscovery(entityEvents notifiers.EventWriter, nodeName string, k8sClient k8s.K8sClient) (*serviceDiscovery, error) {
	if nodeName == "" {
		log.Error().Msg("Node name not set, could not create Kubernetes service discovery")
		return nil, fmt.Errorf("node name not set")
	}

	kd := &serviceDiscovery{
		cli:          k8sClient.GetClientSet(),
		nodeName:     nodeName,
		entityEvents: entityEvents,
	}
	return kd, nil
}

func (kd *serviceDiscovery) start() {
	kd.ctx, kd.cancel = context.WithCancel(context.Background())

	kd.waitGroup.Add(1)
	panichandler.Go(func() {
		defer kd.waitGroup.Done()

		stopChan := make(chan struct{})
		defer close(stopChan)

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

			informerFactory := informers.NewSharedInformerFactoryWithOptions(kd.cli, time.Second*5, informers.WithNamespace(metav1.NamespaceAll))
			endpointsInformer := informerFactory.Core().V1().Endpoints().Informer()
			_, err = endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					endpoints := obj.(*v1.Endpoints)
					kd.updateEntitiesFromEndpoints(endpoints, add)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					oldEndpoints := oldObj.(*v1.Endpoints)
					newEndpoints := newObj.(*v1.Endpoints)
					if newEndpoints.ResourceVersion == oldEndpoints.ResourceVersion {
						return
					}
					endpointsItem, ok, ierr := endpointsInformer.GetIndexer().GetByKey(oldEndpoints.Namespace + "/" + oldEndpoints.Name)
					if err != nil {
						log.Error().Err(ierr).Msg("Could not get endpoints from cache")
						return
					}
					if ok {
						oldEndpointsFromCache := endpointsItem.(*v1.Endpoints)
						kd.updateEntitiesFromEndpoints(oldEndpointsFromCache, remove)
					}
					kd.updateEntitiesFromEndpoints(newEndpoints, add)
				},
				DeleteFunc: func(obj interface{}) {
					endpoints := obj.(*v1.Endpoints)
					kd.updateEntitiesFromEndpoints(endpoints, remove)
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Unable to watch endpoints")
				return err
			}

			informerFactory.Start(stopChan)
			synced := informerFactory.WaitForCacheSync(stopChan)
			for v, ok := range synced {
				if !ok {
					err := fmt.Errorf("informer for %s failed to sync", v)
					log.Error().Err(err)
					return err
				}
			}

			<-stopChan

			return nil
		}
		boff := backoff.NewConstantBackOff(5 * time.Second)
		err := backoff.Retry(operation, backoff.WithContext(boff, kd.ctx))
		if err != nil {
			close(stopChan)
		}

		log.Info().Msg("Stopping kubernetes service watcher")
	})
}

func (kd *serviceDiscovery) stop() {
	kd.cancel()
	kd.waitGroup.Wait()
	kd.entityEvents.Purge("")
}

// updatePodInTracker retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) writeEntityInTracker(podInfo podInfo) error {
	entity := &entitycachev1.Entity{
		Services:  podInfo.Services,
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

func (kd *serviceDiscovery) updateEntitiesFromEndpoints(endpoints *v1.Endpoints, operation serviceCacheOperation) {
	pods := getEndpointPods(endpoints, kd.nodeName)

	// get service from endpoints
	svc, err := kd.cli.CoreV1().Services(endpoints.Namespace).Get(kd.ctx, endpoints.Name, metav1.GetOptions{ResourceVersion: endpoints.ResourceVersion})
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get service %s/%s", endpoints.Namespace, endpoints.Name)
	}

	for _, pod := range pods {
		if svc.Spec.ClusterIP != "" {
			if svc.Spec.ClusterIP != "None" {
				pod.ClusterIP = svc.Spec.ClusterIP
			}
		}

		pod.Services = []string{kd.getFQDN(endpoints)}

		if operation == add {
			err := kd.writeEntityInTracker(pod)
			if err != nil {
				log.Error().Msgf("Tracker could not be updated: %v", err)
			}
		} else {
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

// getEndpointPods retrieves a list of pods handled by a given endpoints that are located on a given node.
func getEndpointPods(endpoints *v1.Endpoints, nodeName string) []podInfo {
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
