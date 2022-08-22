package watcher

import (
	"context"
	"errors"
	"path"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
)

// watcher holds the state of the watcher.
type watcher struct {
	waitGroup  sync.WaitGroup
	ctx        context.Context
	trackers   notifiers.Trackers
	etcdClient *etcdclient.Client
	cancel     context.CancelFunc
	etcdPath   string
	revision   int64
}

// Make sure Watcher implements notifiers.Watcher interface.
var _ notifiers.Watcher = (*watcher)(nil)

// NewWatcher creates a new watcher.
func NewWatcher(etcdClient *etcdclient.Client, etcdPath string) (notifiers.Watcher, error) {
	if etcdPath == "" {
		err := errors.New("unable to create etcd watcher because etcdPath is empty")
		log.Error().Err(err).Msg("")
		return nil, err
	}

	w := &watcher{
		etcdPath:   etcdPath,
		etcdClient: etcdClient,
		trackers:   notifiers.NewDefaultTrackers(),
	}

	// context to track the lifecycle of watcher
	// this context gets canceled in stop
	w.ctx, w.cancel = context.WithCancel(context.Background())

	return w, nil
}

// Start first starts trackers, bootstraps from existing keys in etcd and then starts watching etcd prefix w.etcdPath.
func (w *watcher) Start() error {
	err := w.trackers.Start()
	if err != nil {
		return err
	}
	w.bootstrap()

	w.waitGroup.Add(1)

	panichandler.Go(func() {
		defer w.waitGroup.Done()
		// start watch to accumulate events
		// need to start all over again on non-recoverable error in watch response (refer https://pkg.go.dev/github.com/coreos/etcd/clientv3#Watcher)
		wCh := w.etcdClient.Watcher.Watch(clientv3.WithRequireLeader(w.ctx),
			w.etcdPath, clientv3.WithRev(w.revision+1), clientv3.WithPrefix())

		for {
			select {
			case resp, ok := <-wCh:
				if !ok {
					return
				}
				if resp.Canceled {
					log.Error().Err(resp.Err()).Msg("Etcd watch channel was canceled")
					w.trackers.Purge("")
					w.bootstrap()
					continue
				}
				for _, ev := range resp.Events {
					key := getNotifierKey(ev.Kv.Key)
					// Track only the children, skip etcdPath itself
					if path.Clean(string(ev.Kv.Key)) == path.Clean(w.etcdPath) {
						continue
					}

					switch ev.Type {
					case mvccpb.PUT:
						w.trackers.WriteEvent(key, ev.Kv.Value)
					case mvccpb.DELETE:
						w.trackers.RemoveEvent(key)
					}
				}
			case <-w.ctx.Done():
				return
			}
		}
	})
	return nil
}

// Stop cancels watcher context and stops trackers.
func (w *watcher) Stop() error {
	w.cancel()
	w.waitGroup.Wait()
	return w.trackers.Stop()
}

// bootstrap iterates throughout all existing keys in etcd and updates trackers in the existing watcher.
func (w *watcher) bootstrap() {
	var getResp *clientv3.GetResponse
	var err error

	operation := func() error {
		ctx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
		defer cancel()
		getResp, err = w.etcdClient.KV.Get(ctx, w.etcdPath, clientv3.WithPrefix())
		if err != nil {
			log.Error().Err(err).Str("etcdPath", w.etcdPath).Msg("Failed to list keys")
			return err
		}
		return nil
	}

	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = time.Duration(0) // never stop retrying
	boff.MaxInterval = 10 * time.Second

	err = backoff.Retry(operation, backoff.WithContext(boff, w.ctx))
	if err != nil {
		log.Info().Msg("Stopping bootstrap")
		return
	}

	w.revision = getResp.Header.Revision

	kvs := make([]*mvccpb.KeyValue, 0, len(getResp.Kvs))
	kvs = append(kvs, getResp.Kvs...)

	for _, kv := range kvs {
		if string(kv.Key) == w.etcdPath {
			continue
		}
		key := getNotifierKey(kv.Key)
		w.trackers.WriteEvent(key, kv.Value)
	}
}

// AddPrefixNotifier is a helper function to add a new prefix notifier to watcher.
func (w *watcher) AddPrefixNotifier(notifier notifiers.PrefixNotifier) error {
	return w.trackers.AddPrefixNotifier(notifier)
}

// RemovePrefixNotifier is a helper function to remove an existing prefix notifier from watcher.
func (w *watcher) RemovePrefixNotifier(notifier notifiers.PrefixNotifier) error {
	return w.trackers.RemovePrefixNotifier(notifier)
}

// AddKeyNotifier is a helper method to add a new key notifier to watcher.
func (w *watcher) AddKeyNotifier(notifier notifiers.KeyNotifier) error {
	return w.trackers.AddKeyNotifier(notifier)
}

// RemoveKeyNotifier is a helper method to remove an existing key notifier from watcher.
func (w *watcher) RemoveKeyNotifier(notifier notifiers.KeyNotifier) error {
	return w.trackers.RemoveKeyNotifier(notifier)
}

func getNotifierKey(etcdKey []byte) notifiers.Key {
	_, lastElem := path.Split(string(etcdKey))
	return notifiers.Key(lastElem)
}
