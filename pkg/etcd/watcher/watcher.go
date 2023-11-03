package watcher

import (
	"context"
	"errors"
	"path"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
)

// watcher holds the state of the watcher.
type watcher struct {
	notifiers.Trackers
	ctx        context.Context
	etcdClient *etcdclient.Client
	cancel     context.CancelFunc
	etcdPath   string
	waitGroup  panichandler.WaitGroup
}

// Make sure Watcher implements notifiers.Watcher interface.
var _ notifiers.Watcher = (*watcher)(nil)

// NewWatcher creates a new watcher.
func NewWatcher(etcdClient *etcdclient.Client, etcdPath string) (notifiers.Watcher, error) {
	if etcdPath == "" {
		err := errors.New("unable to create etcd watcher because etcdPath is empty")
		log.Error().Err(err).Send()
		return nil, err
	}

	w := &watcher{
		etcdPath:   etcdPath,
		etcdClient: etcdClient,
		Trackers:   notifiers.NewDefaultTrackers(),
	}

	// context to track the lifecycle of watcher
	// this context gets canceled in stop
	w.ctx, w.cancel = context.WithCancel(context.Background())

	return w, nil
}

// Start starts the etcd watcher.
func (w *watcher) Start() error {
	err := w.Trackers.Start()
	if err != nil {
		return err
	}

	w.waitGroup.Go(func() {
		for {
			err := w.doWatch()
			if w.ctx.Err() != nil {
				log.Info().Err(w.ctx.Err()).Msg("Context canceled, stopping etcd watcher")
				return
			}
			log.Error().Err(err).Msg("etcd watch channel was canceled. Re-bootrapping")
		}
	})
	return nil
}

// Stop cancels watcher context and stops trackers.
func (w *watcher) Stop() error {
	w.cancel()
	w.waitGroup.Wait()
	return w.Trackers.Stop()
}

// doWatch iterates throughout all existing keys in etcd, updates trackers in
// the existing watcher, and starts watching for updates.
//
// Errors when context is canceled or watching fails.
func (w *watcher) doWatch() error {
	// bootstrap

	var getResp *clientv3.GetResponse
	operation := func() error {
		ctx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
		defer cancel()
		var err error
		getResp, err = w.etcdClient.Get(ctx, w.etcdPath, clientv3.WithPrefix())
		if err != nil {
			log.Error().Err(err).Str("etcdPath", w.etcdPath).Msg("Failed to list keys")
			return err
		}
		return nil
	}

	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = time.Duration(0) // never stop retrying
	boff.MaxInterval = 10 * time.Second

	err := backoff.Retry(operation, backoff.WithContext(boff, w.ctx))
	if err != nil {
		log.Info().Msg("Stopping bootstrap")
		return err
	}

	// Remove all potentially stale keys from the Trackers, before writing new state.
	w.Trackers.Purge("")

	for _, kv := range getResp.Kvs {
		if string(kv.Key) == w.etcdPath {
			continue
		}
		key := getNotifierKey(kv.Key)
		w.WriteEvent(key, kv.Value)
	}

	// start watch to accumulate events
	// Note: Watching from "Get" revision + 1, because revision passed to
	// WithRev is inclusive and we want to watch for all events occurring
	// _after_ we retrieved the initial state.
	wCh, err := w.etcdClient.Watch(clientv3.WithRequireLeader(w.ctx),
		w.etcdPath, clientv3.WithRev(getResp.Header.Revision+1), clientv3.WithPrefix())
	if err != nil {
		log.Error().Err(err).Msg("Failed to start etcd watch")
		return err
	}

	for resp := range wCh {
		if err := resp.Err(); err != nil {
			// need to start all over again on non-recoverable error in watch
			// response (refer https://pkg.go.dev/github.com/coreos/etcd/clientv3#Watcher)
			return err
		}
		for _, ev := range resp.Events {
			key := getNotifierKey(ev.Kv.Key)
			// Track only the children, skip etcdPath itself
			if path.Clean(string(ev.Kv.Key)) == path.Clean(w.etcdPath) {
				continue
			}

			switch ev.Type {
			case mvccpb.PUT:
				w.WriteEvent(key, ev.Kv.Value)
			case mvccpb.DELETE:
				w.RemoveEvent(key)
			}
		}
	}

	// This probably should not happen, as watcher should send Err() != nil
	// response response on the wCh before closing it.
	return nil
}

func getNotifierKey(etcdKey []byte) notifiers.Key {
	_, lastElem := path.Split(string(etcdKey))
	return notifiers.Key(lastElem)
}
