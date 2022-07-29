package notifier

import (
	clientv3 "go.etcd.io/etcd/client/v3"

	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// PrefixToEtcdNotifier holds the state of a notifier that writes raw/transformed contents of a watched prefix to etcd.
type PrefixToEtcdNotifier struct {
	notifiers.PrefixNotifierBase
	etcdClient *etcdclient.Client
	etcdWriter *etcdwriter.Writer
	etcdPath   string
}

// Make sure PrefixToEtcdNotifier implements PrefixNotifier.
var _ notifiers.PrefixNotifier = (*PrefixToEtcdNotifier)(nil)

// NewPrefixToEtcdNotifier returns a new prefix notifier that writes raw/transformed contents to etcd at "etcdPath/key".
func NewPrefixToEtcdNotifier(
	etcdPath string,
	etcdClient *etcdclient.Client,
	withLease bool,
) *PrefixToEtcdNotifier {
	pen := &PrefixToEtcdNotifier{
		etcdPath:   etcdPath,
		etcdClient: etcdClient,
		etcdWriter: etcdwriter.NewWriter(etcdClient, withLease),
	}
	return pen
}

// Start starts the prefix notifier.
func (pen *PrefixToEtcdNotifier) Start() error {
	// purge etcd path -- as OnStart hooks are executed in order, this would be the first operation on the writer
	pen.etcdWriter.Delete(pen.etcdPath, clientv3.WithPrefix())
	return nil
}

// Stop stops the prefix notifier.
func (pen *PrefixToEtcdNotifier) Stop() error {
	pen.etcdWriter.Close()
	return nil
}

// GetKeyNotifier gets the underlying key notifier from prefix notifier.
func (pen *PrefixToEtcdNotifier) GetKeyNotifier(key notifiers.Key) notifiers.KeyNotifier {
	ken := &KeyToEtcdNotifier{
		etcdWriter: pen.etcdWriter,
		etcdPath:   pen.etcdPath,
	}
	return ken
}
