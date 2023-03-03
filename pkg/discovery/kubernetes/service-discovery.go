package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sourcegraph/conc/stream"
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

type podServiceUpdate struct {
	Name      string
	Namespace string
	NodeName  string
	UID       string
	IPAddress string
	Service   string
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
	kd.updateEndpoints(endpoints)
}

func (kd *serviceDiscovery) handleEndpointsUpdate(oldObj, newObj interface{}) {
	oldEndpoints := oldObj.(*v1.Endpoints)
	newEndpoints := newObj.(*v1.Endpoints)
	if oldEndpoints.ResourceVersion == newEndpoints.ResourceVersion {
		return
	}
	// make a deep copy of oldEndpoints
	toRemove := oldEndpoints.DeepCopy()
	// from this copy remove addresses that are present in newEndpoints
	for _, newSubsets := range newEndpoints.Subsets {
		for _, newAddresses := range newSubsets.Addresses {
			for i, oldSubsets := range toRemove.Subsets {
				for j, oldAddresses := range oldSubsets.Addresses {
					if newAddresses.TargetRef.UID == oldAddresses.TargetRef.UID {
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

func (kd *serviceDiscovery) getEntityFromTracker(uid string) *entitycachev1.Entity {
	bytes := kd.entityEvents.GetCurrentValue(notifiers.Key(uid))
	if bytes == nil {
		return nil
	}
	entity := &entitycachev1.Entity{}
	err := json.Unmarshal(bytes, entity)
	if err != nil {
		log.Error().Err(err).Msg("Could not unmarshal entity")
		return nil
	}
	return entity
}

// getService return the full qualified domain name of a given service.
func (kd *serviceDiscovery) getService(endpoints *v1.Endpoints) string {
	name := endpoints.Name
	namespace := endpoints.Namespace

	service := fmt.Sprintf("%s.%s.svc.%s", name, namespace, kd.clusterDomain)
	return service
}

func (kd *serviceDiscovery) removeEndpoints(endpoints *v1.Endpoints) {
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			entity := kd.getEntityFromTracker(string(address.TargetRef.UID))
			if entity != nil {
				entity.Services = utils.RemoveFromSlice(entity.Services, kd.getService(endpoints))
				if kd.shouldRemove(entity) {
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

// updatePodService retrieves stored pod data from tracker, enriches it with new info and send the updated version.
func (kd *serviceDiscovery) updatePodService(pod podServiceUpdate) error {
	entity := kd.getEntityFromTracker(pod.UID)
	if entity != nil {
		// append to services if it doesn't exist
		if !utils.SliceContains(entity.Services, pod.Service) {
			entity.Services = append(entity.Services, pod.Service)
		}
	} else {
		// create new entity
		entity = &entitycachev1.Entity{
			Services:  []string{pod.Service},
			IpAddress: pod.IPAddress,
			Namespace: pod.Namespace,
			NodeName:  pod.NodeName,
			Uid:       pod.UID,
			Prefix:    podTrackerPrefix,
			Name:      pod.Name,
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

func (kd *serviceDiscovery) shouldRemove(entity *entitycachev1.Entity) bool {
	// once we have more informers then add additional checks here
	return len(entity.Services) == 0
}

func (kd *serviceDiscovery) updateEndpoints(endpoints *v1.Endpoints) {
	service := kd.getService(endpoints)
	pods := kd.getServicePods(endpoints)

	kd.serviceStream.Go(func() stream.Callback {
		return kd.fetchService(endpoints, pods, service)
	})
}

func (kd *serviceDiscovery) fetchService(endpoints *v1.Endpoints, pods []podServiceUpdate, service string) stream.Callback {
	// clusterIP := ""
	// op := func() error {
	// 	svc, err := kd.cli.CoreV1().Services(endpoints.Namespace).Get(kd.ctx, endpoints.Name, metav1.GetOptions{})
	// 	if kd.ctx.Err() != nil {
	// 		return backoff.Permanent(kd.ctx.Err())
	// 	}
	// 	if err != nil {
	// 		log.Trace().Err(err).Str("namespace", endpoints.Namespace).Str("name", endpoints.Name).Msg("Could not fetch service")
	// 		return nil
	// 	}
	// 	if svc.Spec.Type != v1.ServiceTypeClusterIP {
	// 		return nil
	// 	}
	// 	clusterIP = svc.Spec.ClusterIP
	// 	return nil
	// }
	// err := backoff.Retry(op, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
	// if err != nil {
	// 	log.Error().Err(err).Str("namespace", endpoints.Namespace).Str("name", endpoints.Name).Msg("Context canceled while fetching service")
	// 	return func() {}
	// }

	return func() {
		for _, pod := range pods {
			p := pod
			err := kd.updatePodService(p)
			if err != nil {
				log.Error().Err(err).Msg("Tracker could not be updated")
			}

		}
	}
}

// getServicePods retrieves a list of pods handled by a given service that are located on a given node.
func (kd *serviceDiscovery) getServicePods(endpoints *v1.Endpoints) []podServiceUpdate {
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
				Service:   kd.getService(endpoints),
			}
			if address.NodeName != nil {
				p.NodeName = *address.NodeName
			}
			pods = append(pods, p)
		}
	}

	return pods
}
