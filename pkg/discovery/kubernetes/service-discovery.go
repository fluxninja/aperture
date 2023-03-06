package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sourcegraph/conc/stream"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/utils"
)

const podTrackerPrefix = "kubernetes_pod"

type podServiceUpdate struct {
	Name      string
	Namespace string
	NodeName  string
	UID       string
	IPAddress string
	Service   string
	Labels    map[string]string
}

// serviceDiscovery is a collector that collects Kubernetes information periodically.
type serviceDiscovery struct {
	ctx             context.Context
	cli             kubernetes.Interface
	informerFactory informers.SharedInformerFactory
	entityEvents    notifiers.EventWriter
	cancel          context.CancelFunc
	serviceStream   *stream.Stream
	clusterDomain   string
}

func newServiceDiscovery(
	entityEvents notifiers.EventWriter,
	k8sClient k8s.K8sClient,
) (*serviceDiscovery, error) {
	kd := &serviceDiscovery{
		cli:           k8sClient.GetClientSet(),
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

			serviceInformer := kd.informerFactory.Core().V1().Services().Informer()
			_, err = serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc:    kd.handleServiceAdd,
				UpdateFunc: kd.handleServiceUpdate,
				DeleteFunc: kd.handleServiceDelete,
			})
			if err != nil {
				return err
			}

			kd.informerFactory.Start(kd.ctx.Done())
			// wait for the caches to be synced
			if !cache.WaitForCacheSync(kd.ctx.Done(), endpointsInformer.HasSynced, serviceInformer.HasSynced) {
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

// Endpoints informer handlers

func (kd *serviceDiscovery) handleEndpointsAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.updateEndpoints(endpoints)
}

func (kd *serviceDiscovery) handleEndpointsUpdate(oldObj, newObj interface{}) {
	oldEndpoints := oldObj.(*v1.Endpoints)
	newEndpoints := newObj.(*v1.Endpoints)
	// make a deep copy of oldEndpoints
	toRemove := oldEndpoints.DeepCopy()
	// check if an address in toRemove is found in any subsets of newEndpoints, if so remove it from toRemove
	for _, newSubset := range newEndpoints.Subsets {
		for _, newAddress := range newSubset.Addresses {
			for i, oldSubset := range toRemove.Subsets {
				for j, oldAddress := range oldSubset.Addresses {
					if newAddress.TargetRef.UID == oldAddress.TargetRef.UID {
						toRemove.Subsets[i].Addresses = append(toRemove.Subsets[i].Addresses[:j], toRemove.Subsets[i].Addresses[j+1:]...)
					}
				}
			}
		}
	}
	kd.removeEndpoints(toRemove)
	kd.updateEndpoints(newEndpoints)
}

func (kd *serviceDiscovery) handleEndpointsDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	kd.removeEndpoints(endpoints)
}

func (kd *serviceDiscovery) updateEndpoints(endpoints *v1.Endpoints) {
	pods := kd.getPodServiceUpdatesFromEndpoints(endpoints)
	for _, pod := range pods {
		p := pod
		err := kd.updatePodServiceEntity(p)
		if err != nil {
			log.Error().Err(err).Msg("Tracker could not be updated")
		}
	}
}

// getPodServiceUpdatesFromEndpoints retrieves a list of pods handled by a given service that are located on a given node.
func (kd *serviceDiscovery) getPodServiceUpdatesFromEndpoints(endpoints *v1.Endpoints) []podServiceUpdate {
	var pods []podServiceUpdate
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			if address.TargetRef == nil {
				continue
			}
			if address.TargetRef.Kind != "Pod" {
				continue
			}
			p := podServiceUpdate{
				Name:      address.TargetRef.Name,
				Namespace: address.TargetRef.Namespace,
				UID:       string(address.TargetRef.UID),
				IPAddress: address.IP,
				Service:   kd.getService(endpoints.Namespace, endpoints.Name),
			}
			if address.NodeName != nil {
				p.NodeName = *address.NodeName
			}
			pods = append(pods, p)
		}
	}
	return pods
}

func (kd *serviceDiscovery) removeEndpoints(endpoints *v1.Endpoints) {
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			entity := kd.getEntityFromTracker(string(address.TargetRef.UID))
			if entity != nil {
				entity.Services = utils.RemoveFromSlice(entity.Services, kd.getService(endpoints.Namespace, endpoints.Name))
				if shouldRemove(entity) {
					kd.entityEvents.RemoveEvent(notifiers.Key(address.TargetRef.UID))
				} else {
					bytes, err := json.Marshal(entity)
					if err != nil {
						log.Error().Err(err).Msg("Could not marshal entity")
						continue
					}
					kd.entityEvents.WriteEvent(notifiers.Key(address.TargetRef.UID), bytes)
				}
			}
		}
	}
}

// Service informer handlers

func (kd *serviceDiscovery) handleServiceAdd(obj interface{}) {
	service := obj.(*v1.Service)
	kd.addClusterIPs(service)
}

func (kd *serviceDiscovery) handleServiceUpdate(oldObj, newObj interface{}) {
	oldService := oldObj.(*v1.Service)
	newService := newObj.(*v1.Service)
	// make a deep copy of oldService
	toRemove := oldService.DeepCopy()
	// from this copy remove clusterIPs that are present in newService
	for _, newClusterIP := range newService.Spec.ClusterIPs {
		for i, oldClusterIP := range toRemove.Spec.ClusterIPs {
			if newClusterIP == oldClusterIP {
				toRemove.Spec.ClusterIPs = append(toRemove.Spec.ClusterIPs[:i], toRemove.Spec.ClusterIPs[i+1:]...)
			}
		}
	}
	kd.removeClusterIPs(toRemove)
	kd.addClusterIPs(newService)
}

func (kd *serviceDiscovery) handleServiceDelete(obj interface{}) {
	service := obj.(*v1.Service)
	kd.removeClusterIPs(service)
}

func (kd *serviceDiscovery) addClusterIPs(service *v1.Service) {
	clusterIPs := getClusterIPsFromService(service)
	for _, clusterIP := range clusterIPs {
		p := podServiceUpdate{
			Name:      strings.Join([]string{"ClusterIP", service.Namespace, service.Name, clusterIP}, "-"),
			Namespace: service.Namespace,
			UID:       clusterIP,
			IPAddress: clusterIP,
			Service:   kd.getService(service.Namespace, service.Name),
		}
		err := kd.updatePodServiceEntity(p)
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

func (kd *serviceDiscovery) getEntityFromTracker(uid string) *entitiesv1.Entity {
	bytes := kd.entityEvents.GetCurrentValue(notifiers.Key(uid))
	if bytes == nil {
		return nil
	}
	entity := &entitiesv1.Entity{}
	err := json.Unmarshal(bytes, entity)
	if err != nil {
		log.Error().Err(err).Msg("Could not unmarshal entity")
		return nil
	}
	return entity
}

// getService return the full qualified domain name of a given service.
func (kd *serviceDiscovery) getService(namespace, name string) string {
	service := fmt.Sprintf("%s.%s.svc.%s", name, namespace, kd.clusterDomain)
	return service
}

// updatePodServiceEntity retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) updatePodServiceEntity(pod podServiceUpdate) error {
	entity := kd.getEntityFromTracker(pod.UID)
	if entity != nil {
		// append to services if it doesn't exist
		if pod.Service != "" && !utils.SliceContains(entity.Services, pod.Service) {
			entity.Services = append(entity.Services, pod.Service)
		}
	} else {
		// create new entity
		entity = &entitiesv1.Entity{
			Prefix:    podTrackerPrefix,
			Uid:       pod.UID,
			IpAddress: pod.IPAddress,
			Name:      pod.Name,
			Namespace: pod.Namespace,
			NodeName:  pod.NodeName,
		}
		if pod.Service != "" {
			entity.Services = []string{pod.Service}
		}
	}

	value, err := json.Marshal(entity)
	if err != nil {
		log.Error().Msgf("Error marshaling entity: %v", err)
		return err
	}
	kd.entityEvents.WriteEvent(notifiers.Key(pod.UID), value)
	return nil
}

/* Helper functions */

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
