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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func newControlPointDiscovery(election *election.Election, k8sClient k8s.K8sClient) (*controlPointDiscovery, error) {
	if k8sClient.GetErrNotInCluster() {
		log.Info().Msg("Not in Kubernetes cluster, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}
	if k8sClient.GetErr() != nil {
		log.Error().Err(k8sClient.GetErr()).Msg("Error when creating Kubernetes client, could not create Kubernetes service discovery")
		return nil, k8sClient.GetErr()
	}

	cpd := &controlPointDiscovery{
		cli:         k8sClient.GetClientSet(),
		election:    election,
		cacheStores: make(map[string]cache.Store),
	}

	return cpd, nil
}

// controlPointDiscovery is a struct that helps with Kubernetes control point discovery.
type controlPointDiscovery struct {
	waitGroup   sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	cli         *kubernetes.Clientset
	cacheStores map[string]cache.Store
	election    *election.Election
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
			_, apiResourceListList, err := cpd.cli.DiscoveryClient.ServerGroupsAndResources()
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

			log.Info().Msgf("Scalable resources: %v", scalableResources)

			// Cache Store for each scalable resource
			/*for _, scalableResource := range scalableResources {
			  store := cache.NewStore(cache.MetaNamespaceKeyFunc)
			}*/

			<-cpd.ctx.Done()
			return nil
		}

		boff := backoff.NewConstantBackOff(5 * time.Second)
		_ = backoff.Retry(operation, backoff.WithContext(boff, cpd.ctx))

		log.Info().Msg("Stopping kubernetes control point watcher")
	})
}

func (cpd *controlPointDiscovery) stop() {
	cpd.cancel()
	cpd.waitGroup.Wait()
}
