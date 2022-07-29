package etcd

import (
	"context"
	"runtime"

	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panic"
	"github.com/lukejoshuapark/infchan"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	put = 0
	del = 1
)

type op struct {
	key    string
	value  []byte
	opts   []clientv3.OpOption
	opType int
}

// Writer holds fields for etcd writer.
type Writer struct {
	context    context.Context
	etcdClient *etcdclient.Client
	opChannel  infchan.Channel[op]
	cancel     context.CancelFunc
	withLease  bool
}

// NewWriter returns a new etcd writer.
func NewWriter(etcdClient *etcdclient.Client, withLease bool, opts ...clientv3.OpOption) *Writer {
	ew := &Writer{
		etcdClient: etcdClient,
		withLease:  withLease,
		opChannel:  infchan.NewChannel[op](),
	}
	// Set finalizer to automatically close channel
	runtime.SetFinalizer(ew, func(ew *Writer) {
		// drain the ew.opChannel.Out
		for {
			select {
			case <-ew.opChannel.Out():
			default:
				ew.opChannel.Close()
				return
			}
		}
	})

	ew.context, ew.cancel = context.WithCancel(context.Background())

	panic.Go(func() {
		// start processing ops
		for {
			select {
			case opt := <-ew.opChannel.Out():
				opt.opts = append(opt.opts, opts...)
				if ew.withLease && opt.opType == put {
					opt.opts = append(opt.opts, clientv3.WithLease(ew.etcdClient.LeaseID))
				}

				var err error
				switch opt.opType {
				case put:
					_, err = ew.etcdClient.KV.Put(clientv3.WithRequireLeader(ew.context), opt.key, string(opt.value), opt.opts...)
				case del:
					_, err = ew.etcdClient.KV.Delete(clientv3.WithRequireLeader(ew.context), opt.key, opt.opts...)
				}
				if err != nil {
					log.Error().Err(err).Msg("failed to write to etcd")
				}
			case <-ew.context.Done():
				return
			}
		}
	})

	return ew
}

// Write writes a key value pair to etcd with options.
func (ew *Writer) Write(key string, value []byte, opts ...clientv3.OpOption) {
	ew.opChannel.In() <- op{opType: put, key: key, value: value, opts: opts}
}

// Delete deletes a key from etcd.
func (ew *Writer) Delete(key string, opts ...clientv3.OpOption) {
	ew.opChannel.In() <- op{opType: del, key: key, opts: opts}
}

// Close closes the etcd writer.
func (ew *Writer) Close() error {
	ew.cancel()
	return nil
}
