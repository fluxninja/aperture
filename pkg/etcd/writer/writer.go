package etcd

import (
	"context"
	"runtime"

	"github.com/lukejoshuapark/infchan"
	clientv3 "go.etcd.io/etcd/client/v3"

	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
)

const (
	put       = 0
	del       = 1
	delPrefix = 2
)

type operation struct {
	key    string
	value  []byte
	opts   []clientv3.OpOption
	opType int
}

// Writer holds fields for etcd writer.
type Writer struct {
	context    context.Context
	etcdClient *etcdclient.Client
	opChannel  infchan.Channel[operation]
	cancel     context.CancelFunc
	withLease  bool
}

// NewWriter returns a new etcd writer.
func NewWriter(etcdClient *etcdclient.Client, withLease bool, opts ...clientv3.OpOption) *Writer {
	ew := &Writer{
		etcdClient: etcdClient,
		opChannel:  infchan.NewChannel[operation](),
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

	panichandler.Go(func() {
		// start processing ops
		for {
			select {
			case op := <-ew.opChannel.Out():
				op.opts = append(op.opts, opts...)
				if ew.withLease && op.opType == put {
					op.opts = append(op.opts, clientv3.WithLease(etcdClient.LeaseID))
				}
				var err error
				switch op.opType {
				case put:
					_, err = etcdClient.KV.Put(clientv3.WithRequireLeader(ew.context), op.key, string(op.value), opts...)
				case del:
					_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ew.context), op.key, opts...)
				case delPrefix:
					opOpts := []clientv3.OpOption{clientv3.WithPrefix()}
					opOpts = append(opOpts, opts...)
					_, err = etcdClient.KV.Delete(clientv3.WithRequireLeader(ew.context), op.key, opOpts...)
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
func (ew *Writer) Write(key string, value []byte) {
	ew.opChannel.In() <- operation{opType: put, key: key, value: value}
}

// Delete deletes a key from etcd.
func (ew *Writer) Delete(key string) {
	ew.opChannel.In() <- operation{opType: del, key: key}
}

// DeletePrefix deletes a whole prefix from etcd.
func (ew *Writer) DeletePrefix(key string) {
	ew.opChannel.In() <- operation{opType: delPrefix, key: key}
}

// Close closes the etcd writer.
func (ew *Writer) Close() error {
	ew.cancel()
	return nil
}
