package kubernetes

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/fluxninja/aperture/operator/api"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

// watcher holds the state of the watcher.
type watcher struct {
	waitGroup sync.WaitGroup
	trackers  notifiers.Trackers
	ctx       context.Context
	cancel    context.CancelFunc
}

// Make sure Watcher implements notifiers.Watcher interface.
var _ notifiers.Watcher = &watcher{}

// NewWatcher prepares watcher instance for the Kuberneter Policy.
func NewWatcher() (*watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	watcher := &watcher{
		trackers: notifiers.NewDefaultTrackers(),
		ctx:      ctx,
		cancel:   cancel,
	}

	return watcher, nil
}

// Start starts the watcher go routines and handles Policy Custom resource events from Kubernetes.
func (w *watcher) Start() error {
	err := w.trackers.Start()
	if err != nil {
		return err
	}

	w.waitGroup.Add(1)

	panichandler.Go(func() {
		defer w.waitGroup.Done()
		operation := func() error {
			scheme := runtime.NewScheme()

			utilruntime.Must(clientgoscheme.AddToScheme(scheme))

			utilruntime.Must(api.AddToScheme(scheme))

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             scheme,
				MetricsBindAddress: "0",
				LeaderElection:     false,
				Namespace:          os.Getenv("APERTURE_CONTROLLER_NAMESPACE"),
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to create Kubernetes Reconciler for Policy")
				return nil
			}

			if err = (&PolicyReconciler{
				Client:   mgr.GetClient(),
				Scheme:   scheme,
				Recorder: mgr.GetEventRecorderFor("aperture-policy"),
				Trackers: w.trackers,
			}).SetupWithManager(mgr); err != nil {
				log.Error().Err(err).Msg("Failed to create Kubernetes controller for policy")
				return nil
			}
			return mgr.Start(w.ctx)
		}

		boff := backoff.NewConstantBackOff(5 * time.Second)

		_ = backoff.Retry(operation, backoff.WithContext(boff, w.ctx))
		log.Info().Msg("Stopping kubernetes watcher for Policy")
	})

	return nil
}

// Stop stops the watcher go routines.
func (w *watcher) Stop() error {
	w.cancel()
	w.waitGroup.Wait()
	return w.trackers.Stop()
}

// AddPrefixNotifier is a helper function to add a new directory notifier to watcher.
func (w *watcher) AddPrefixNotifier(notifier notifiers.PrefixNotifier) error {
	return w.trackers.AddPrefixNotifier(notifier)
}

// RemovePrefixNotifier is a helper function to remove an existing directory notifier from watcher.
func (w *watcher) RemovePrefixNotifier(notifier notifiers.PrefixNotifier) error {
	return w.trackers.RemovePrefixNotifier(notifier)
}

// AddKeyNotifier is a helper method to add a new file notifier to watcher.
func (w *watcher) AddKeyNotifier(notifier notifiers.KeyNotifier) error {
	return w.trackers.AddKeyNotifier(notifier)
}

// RemoveKeyNotifier is a helper method to remove an existing file notifier from watcher.
func (w *watcher) RemoveKeyNotifier(notifier notifiers.KeyNotifier) error {
	return w.trackers.RemoveKeyNotifier(notifier)
}
