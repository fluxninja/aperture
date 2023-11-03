package notifier

import (
	"path"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

// KeyToEtcdNotifier holds the state of a notifier that writes raw/transformed contents of a watched key to another key in etcd.
type KeyToEtcdNotifier struct {
	notifiers.KeyBase
	etcdClient *etcdclient.Client
	etcdPath   string
}

// Make sure KeyToEtcdNotifier implements KeyNotifier.
var _ notifiers.KeyNotifier = (*KeyToEtcdNotifier)(nil)

// NewKeyToEtcdNotifier returns a new notifier that writes raw/transformed contents to etcd at "etcdPath/key".
func NewKeyToEtcdNotifier(
	key notifiers.Key,
	etcdPath string,
	etcdClient *etcdclient.Client,
	withLease bool,
) (*KeyToEtcdNotifier, error) {
	return newKeyToEtcdNotifier(key, etcdPath, etcdClient)
}

func newKeyToEtcdNotifier(
	key notifiers.Key,
	etcdPath string,
	etcdClient *etcdclient.Client,
) (*KeyToEtcdNotifier, error) {
	ken := &KeyToEtcdNotifier{
		KeyBase:    notifiers.NewKeyBase(key),
		etcdPath:   etcdPath,
		etcdClient: etcdClient,
	}
	return ken, nil
}

// Start starts the key notifier.
func (ken *KeyToEtcdNotifier) Start() error {
	// delete existing key on start
	ken.etcdClient.Delete(path.Join(ken.etcdPath, ken.GetKey().String()))
	return nil
}

// Stop stops the key notifier.
func (ken *KeyToEtcdNotifier) Stop() error {
	return nil
}

// Notify writes/removes to etcd based on received event.
func (ken *KeyToEtcdNotifier) Notify(event notifiers.Event) {
	// Determine etcd key from event: etcdPath + event.Key
	key := path.Join(ken.etcdPath, event.Key.String())

	switch event.Type {
	case notifiers.Write:
		ken.etcdClient.Put(key, string(event.Value))
	case notifiers.Remove:
		ken.etcdClient.Delete(key)
	}
}
