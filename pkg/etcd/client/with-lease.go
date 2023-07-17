package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// kvWithLease implements KV by attaching Lease to every Put it makes.
//
// This makes all the keys it creates scoped to the session.
type kvWithLease struct {
	lease clientv3.LeaseID
	rawKV clientv3.KV
}

// Get implements clientv3.KV.
func (c kvWithLease) Get(
	ctx context.Context,
	key string,
	opts ...clientv3.OpOption,
) (*clientv3.GetResponse, error) {
	return c.rawKV.Get(ctx, key, opts...)
}

// Put implements clientv3.KV.
func (c kvWithLease) Put(
	ctx context.Context,
	key, val string,
	opts ...clientv3.OpOption,
) (*clientv3.PutResponse, error) {
	return c.rawKV.Put(ctx, key, val, append(opts, clientv3.WithLease(c.lease))...)
}

// Delete implements clientv3.KV.
func (c kvWithLease) Delete(
	ctx context.Context,
	key string,
	opts ...clientv3.OpOption,
) (*clientv3.DeleteResponse, error) {
	return c.rawKV.Delete(ctx, key, opts...)
}

// Compact implements clientv3.KV.
func (c kvWithLease) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return c.rawKV.Compact(ctx, rev, opts...)
}

// Do implements clientv3.KV.
func (c kvWithLease) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	// It's possible to implement, but we don't have a need for it now.
	log.Panic().Msg("kvWithLease.Do unimplemented")
	return clientv3.OpResponse{}, nil
}

// Txn implements clientv3.KV.
func (c kvWithLease) Txn(ctx context.Context) clientv3.Txn {
	// It's possible to implement, but we don't have a need for it now.
	log.Panic().Msg("kvWithLease.Txn unimplemented")
	return nil
}
