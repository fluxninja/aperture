package discovery

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/k8s"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
)

func newControlPointDiscovery(etcdClient *etcdclient.Client, k8sClient k8s.K8sClient, controlPointStore AutoScaleControlPointStore) (*controlPointDiscovery, error) {
	cpd := &controlPointDiscovery{
		etcdClient:        etcdClient,
		controlPointStore: controlPointStore,
		discoveryClient:   k8sClient.GetClientSet().DiscoveryClient,
		dynamicClient:     k8sClient.GetDynamicClient(),
	}

	etcdClient.AddElectionWatcher(cpd)

	return cpd, nil
}

// controlPointDiscovery is a struct that helps with Kubernetes control point discovery.
type controlPointDiscovery struct {
	waitGroup         panichandler.WaitGroup
	ctx               context.Context
	cancel            context.CancelFunc
	controlPointStore AutoScaleControlPointStore
	discoveryClient   discovery.DiscoveryInterface
	dynamicClient     dynamic.Interface
	etcdClient        *etcdclient.Client
}

// controlPointDiscovery implements the etcdclient.ElectionWatcher interface.
var _ etcdclient.ElectionWatcher = (*controlPointDiscovery)(nil)

// OnLeaderStart is called when this instance becomes the leader.
func (cpd *controlPointDiscovery) OnLeaderStart() {
	log.Info().Msg("Starting kubernetes control point discovery")
	cpd.start()
}

// OnLeaderStop is called when this instance stops being the leader.
func (cpd *controlPointDiscovery) OnLeaderStop() {
	log.Info().Msg("Stopping kubernetes control point discovery")
	cpd.stop()
}

// Start starts the Kubernetes control point discovery.
func (cpd *controlPointDiscovery) start() {
	cpd.ctx, cpd.cancel = context.WithCancel(context.Background())

	cpd.waitGroup.Go(func() {
		operation := func() error {
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

			groupVersionResourceSet := make(map[schema.GroupVersionResource]string)
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
						groupVersionResourceSet[groupVersionResource] = apiResource.Kind
					}
				}
			}

			// Track each scalable resource
			for groupVersionResource, kind := range groupVersionResourceSet {
				log.Info().Msgf("Starting watch for Group: %s, Version: %s, Resource: %s", groupVersionResource.Group, groupVersionResource.Version, groupVersionResource.Resource)
				resourceInterface := cpd.dynamicClient.Resource(groupVersionResource)

				// watch for changes
				_, controller := cache.NewInformer(
					&cache.ListWatch{
						ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
							return resourceInterface.List(cpd.ctx, options)
						},
						WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
							return resourceInterface.Watch(cpd.ctx, options)
						},
					},
					&unstructured.Unstructured{},
					0,
					cpd.createResourceEventHandlerFuncs(groupVersionResource, kind),
				)

				resourceWatcher := resourceWatcher{
					controller:           controller,
					groupVersionResource: groupVersionResource,
					ctx:                  cpd.ctx,
				}

				resourceWatcher.goRun()
			}

			return nil
		}

		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, cpd.ctx))
	})
}

func (cpd *controlPointDiscovery) createResourceEventHandlerFuncs(groupVersionResource schema.GroupVersionResource, kind string) cache.ResourceEventHandlerFuncs {
	controlPointFromObject := func(obj interface{}) AutoScaleControlPoint {
		// read the name of the resource
		name := obj.(*unstructured.Unstructured).GetName()
		namespace := obj.(*unstructured.Unstructured).GetNamespace()
		return AutoScaleControlPoint{
			Group:     groupVersionResource.Group,
			Version:   groupVersionResource.Version,
			Kind:      kind,
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

type resourceWatcher struct {
	controller           cache.Controller
	ctx                  context.Context
	groupVersionResource schema.GroupVersionResource
}

func (rw *resourceWatcher) goRun() {
	panichandler.Go(func() {
		// Run controller
		rw.controller.Run(rw.ctx.Done())
		log.Info().Msg("Stopped kubernetes control point watcher for resource " + rw.groupVersionResource.String())
	})
}
