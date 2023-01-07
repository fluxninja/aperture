package kubernetes

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func newControlPointDiscovery(election *election.Election, k8sClient k8s.K8sClient, discoveryClient discovery.DiscoveryInterface, dynClient dynamic.Interface, controlPointStore ControlPointStore) (*controlPointDiscovery, error) {
	if k8sClient.GetErrNotInCluster() {
		log.Info().Msg("Not in Kubernetes cluster, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}
	if k8sClient.GetErr() != nil {
		log.Error().Err(k8sClient.GetErr()).Msg("Error when creating Kubernetes client, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}

	cpd := &controlPointDiscovery{
		cli:               k8sClient.GetClientSet(),
		election:          election,
		controlPointStore: controlPointStore,
		discoveryClient:   discoveryClient,
		dynClient:         dynClient,
	}

	return cpd, nil
}

// controlPointDiscovery is a struct that helps with Kubernetes control point discovery.
type controlPointDiscovery struct {
	waitGroup         sync.WaitGroup
	ctx               context.Context
	cancel            context.CancelFunc
	cli               *kubernetes.Clientset
	controlPointStore ControlPointStore
	discoveryClient   discovery.DiscoveryInterface
	dynClient         dynamic.Interface
	election          *election.Election
}

// Start starts the Kubernetes control point discovery.
func (cpd *controlPointDiscovery) start() {
	cpd.ctx, cpd.cancel = context.WithCancel(context.Background())

	cpd.waitGroup.Add(1)

	panichandler.Go(func() {
		defer cpd.waitGroup.Done()

		operation := func() error {
			// Proceed only if we are the leader
			for {
				if cpd.election.IsLeader() {
					// Proceed
					break
				} else {
					// Check again in 5 seconds
					time.Sleep(5 * time.Second)
				}
			}

			// Discover all resources with /scale subresource
			_, apiResourceListList, err := cpd.discoveryClient.ServerGroupsAndResources()
			if err != nil {
				log.Error().Err(err).Msg("Unable to get API resource list")
				return err
			}
			scalableResources := []string{}

			for _, apiResourceList := range apiResourceListList {
				for _, apiResource := range apiResourceList.APIResources {
					if apiResource.Kind == "Scale" {
						// Get its parent resource
						parentResource := strings.TrimSuffix(apiResource.Name, "/scale")
						scalableResources = append(scalableResources, parentResource)
					}
				}
			}

			groupVersionResourceSet := make(map[schema.GroupVersionResource]interface{})
			log.Info().Msgf("Scalable resources: %v", scalableResources)
			for _, apiResourceList := range apiResourceListList {
				for _, apiResource := range apiResourceList.APIResources {
					// Check if apiResource.Name belongs to scalableResources
					if slices.Contains(scalableResources, apiResource.Name) {
						groupVersion, parseErr := schema.ParseGroupVersion(apiResourceList.GroupVersion)
						if parseErr != nil {
							log.Error().Err(parseErr).Msg("Unable to parse group version")
							return parseErr
						}
						groupVersionResource := schema.GroupVersionResource{
							Group:    groupVersion.Group,
							Version:  groupVersion.Version,
							Resource: apiResource.Name,
						}
						// Add to groupVersionResourceSet
						groupVersionResourceSet[groupVersionResource] = nil
					}
				}
			}

			// Track each scalable resource
			for groupVersionResource := range groupVersionResourceSet {
				log.Info().Msgf("Starting watch for Group: %s, Version: %s, Resource: %s", groupVersionResource.Group, groupVersionResource.Version, groupVersionResource.Resource)
				resourceInterface := cpd.dynClient.Resource(groupVersionResource)

				// watch for changes
				_, controller := cache.NewInformer(
					&cache.ListWatch{
						ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
							return resourceInterface.List(context.Background(), options)
						},
						WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
							return resourceInterface.Watch(context.Background(), options)
						},
					},
					&unstructured.Unstructured{},
					0,
					cpd.createResourceEventHandlerFuncs(groupVersionResource),
				)

				// start controller
				panichandler.Go(func() {
					controller.Run(cpd.ctx.Done())
				})
			}

			<-cpd.ctx.Done()
			return nil
		}

		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, cpd.ctx))

		log.Info().Msg("Stopped kubernetes control point watcher")
	})
}

func (cpd *controlPointDiscovery) createResourceEventHandlerFuncs(groupVersionResource schema.GroupVersionResource) cache.ResourceEventHandlerFuncs {
	controlPointFromObject := func(obj interface{}) ControlPoint {
		// read the name of the resource
		name := obj.(*unstructured.Unstructured).GetName()
		namespace := obj.(*unstructured.Unstructured).GetNamespace()
		return ControlPoint{
			Group:     groupVersionResource.Group,
			Version:   groupVersionResource.Version,
			Type:      groupVersionResource.Resource,
			Name:      name,
			Namespace: namespace,
		}
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			controlPoint := controlPointFromObject(obj)
			cpd.controlPointStore.Add(controlPoint)
		},
		UpdateFunc: func(_, obj interface{}) {
			controlPoint := controlPointFromObject(obj)
			cpd.controlPointStore.Update(controlPoint)
		},
		DeleteFunc: func(obj interface{}) {
			controlPoint := controlPointFromObject(obj)
			cpd.controlPointStore.Delete(controlPoint)
		},
	}
}

func (cpd *controlPointDiscovery) stop() {
	cpd.cancel()
	cpd.waitGroup.Wait()
}
